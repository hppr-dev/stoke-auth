arguments_test() {
	DESCRIPTION="Run tests"
	SUBCOMMANDS="int|unit|client|clean"
	INT_DESCRIPTION="Run integration tests"
	INT_OPTIONS="image:i:str rimage:R:str build:b:str case:C:str log:l:bool logfile:L:bool progress:P:bool all:a:bool cert:c:bool db:d:bool race:r:bool provider:p:bool env:e:bool"
	UNIT_DESCRIPTION="Run unit tests"
	UNIT_OPTIONS="cover:c:bool html:h:bool func:f:bool name:n:str"
}

task_test() {
	cd $TASK_DIR
	if [[ "$TASK_SUBCOMMAND" == "int" ]]
	then
		if [[ -z "$ARG_IMAGE" ]]
		then
			ARG_IMAGE=$( docker image ls | grep stoke-inttest | head -n 1 | awk '{print $1}' )
			echo Running with recent found image: $ARG_IMAGE
		fi

		if [[ -z "$ARG_RIMAGE" ]]
		then
			ARG_RIMAGE=$( docker image ls | grep stoke-race | head -n 1 | awk '{print $1}' )
			echo Running with recent found race image: $ARG_RIMAGE
		fi

		if [[ -n "$ARG_BUILD" ]] || [[ -z "$ARG_IMAGE" ]]
		then
			ARG_IMAGE=stoke-inttest-$(date +%Y%m%d%H%M)
			echo Building new image $ARG_IMAGE...
			docker build -t $ARG_IMAGE $TASK_DIR
		fi

		if [[ -n "$ARG_BUILD" ]] || [[ -z "$ARG_RIMAGE" ]]
		then
			ARG_IMAGE=stoke-race-$(date +%Y%m%d%H%M)
			echo Building new image $ARG_IMAGE...
			docker build --build-arg EXTRA_BUILD_ARGS=-race -t $ARG_IMAGE $TASK_DIR
		fi

		echo Starting supplemental containers...
		cd $TASK_DIR/test
		docker compose up -d

		if [[ -z "$ARG_CERT$ARG_DB$ARG_RACE$ARG_PROVIDER$ARG_ENV" ]]
		then
			ARG_ALL="yes"
		fi

		if [[ -n "$ARG_CERT$ARG_ALL" ]]
		then
			# _run_all_configs description config_dir dbinit k6file docker_image pre_start_command
			_run_all_configs "cert smoke" cert_type smoke_test.yaml smoke.js $ARG_IMAGE
		fi


		if [[ -n "$ARG_DB$ARG_ALL" ]]
		then
			_run_all_configs "database smoke" database_type smoke_test.yaml smoke.js $ARG_IMAGE
		fi

		if [[ -n "$ARG_PROVIDER$ARG_ALL" ]]
		then
			_recreate_postgres_schema
			_run_all_configs "provider smoke" provider_type provider_test.yaml provider_test.js $ARG_IMAGE 
		fi

		if [[ -n "$ARG_ENV$ARG_ALL" ]]
		then
			
			_run_all_configs "client/server integration" client_integration client_integration.yaml client_integration.js $ARG_IMAGE _start_client_env
			docker compose -f $TASK_DIR/client/client-test-compose.yaml down
		fi

		if [[ -n "$ARG_RACE$ARG_ALL" ]]
		then
			_run_all_configs "race tests" data_race data_race.yaml data_race.js $ARG_RIMAGE
		fi

		docker compose down

	elif [[ "$TASK_SUBCOMMAND" == "unit" ]]
	then
		echo Running unit tests...
		extra_args=""
		if [[ -n "$ARG_COVER" ]]
		then
			extra_args="-cover"
		fi
		if [[ -n "$ARG_HTML" ]] || [[ -n "$ARG_FUNC" ]]
		then
			extra_args="$extra_args -coverprofile=cover.out"
		fi

		go test $extra_args $ARG_NAME ./internal/{key,usr}

		if [[ -n "$ARG_HTML" ]]
		then
			go tool cover -html cover.out -o cover.html
			rm cover.out
			echo HTML coverage file in $(pwd)/cover.html

		elif [[ -n "$ARG_FUNC" ]]
		then
			go tool cover -func cover.out
			rm cover.out
		fi
	elif [[ "$TASK_SUBCOMMAND" == "clean" ]]
	then
		cd $TASK_DIR/test

		echo Cleaning docker environment...
		docker compose down 
		docker stop stoke-test
		docker rm stoke-test
		docker compose -f $TASK_DIR/client/client-test-compose.yaml down

		echo Removing coverage files...
		if [[ -f "$TASK_DIR/cover.html" ]]
		then
			rm "$TASK_DIR/cover.html"
		fi

		if [[ -f "$TASK_DIR/cover.out" ]]
		then
			rm "$TASK_DIR/cover.out"
		fi

		echo Cleaning test logs file...
		rm -rf $TASK_DIR/test/logs/*
	fi
}

arguments_build() {
	DESCRIPTION="Build the project"
	SUBCOMMANDS="exec|docker"
	DOCKER_REQUIREMENTS="tag:t:str"
}

task_build() {
	cd $TASK_DIR
	if [[ "$TASK_SUBCOMMAND" == "exec" ]]
	then
		echo Building executable...
		cd internal/admin/stoke-admin-ui
		echo Building admin assets...
		npm run build --emptyOutDir

		cd $TASK_DIR
		mkdir build
		go build -o ./build/stoke-server ./cmd/
		echo Executable avilable in $(pwd)/build

	elif [[ "$TASK_SUBCOMMAND" == "docker" ]]
	then
		echo Building image...
		docker build -t hpprdev/stoke-auth:$ARG_TAG .

	fi
}

_run_all_configs() { # desc config_dir dbinit k6file docker_image post_server_command
	echo ====================================================== Running $1 tests...
	for config in ./configs/$2/*
	do
		_run_k6_test $config $3 $4 $5 "$6"
	done


}

_run_k6_test() { # config dbinit k6file docker_image post_server_command
	# Config file relative to CWD
	config=$1
	# dbinit file relative to CWD/configs/dbinit/
	dbinit=$2
	# k6 file relative to CWD/k6/
	k6file=$3
	# docker image
	docker_image=$4
	# command to run after starting the server
	post_server_command="$5"

	if [[ -n "$ARG_CASE" ]] && ! [[ $config =~ "$ARG_CASE" ]]
	then
		echo Did not match case, skipping $config...
		continue
	fi

	docker_args="--name stoke-test \
		-v $(pwd)/$config:/etc/stoke/config.yaml \
		-v $(pwd)/configs/dbinit/$dbinit:/etc/stoke/dbinit.yaml \
		-p 8080:8080 \
		$docker_image -dbinit /etc/stoke/dbinit.yaml"
	docker run -d $docker_args > /dev/null
	if [[ -n "$post_server_command" ]]
	then
		$post_server_command
	fi
	sleep 1


	if ! docker ps | grep stoke-test > /dev/null
	then
		echo Could not start container with $config. 
		docker logs stoke-test
		docker rm stoke-test
		exit
	fi


	if [[ -n "$ARG_LOG" ]]
	then
		echo Tailing container logs...
		docker logs -f stoke-test &
	fi

	if [[ -z "$ARG_PROGRESS" ]]
	then
		k6_args="-q"
	fi

	echo ============================== $config =======================================
	if ! k6 run $k6_args k6/$k6file
	then
		echo "***************************** FAILED ******************************"
	fi

	if [[ -n "$ARG_LOGFILE" ]]
	then
		mkdir -p $TASK_DIR/test/logs
		log_file=$TASK_DIR/test/logs/$(basename $config).json
		echo Copying /etc/stoke/stoke.log file to $log_file
		docker cp stoke-test:/etc/stoke/stoke.log $log_file
	fi

	docker stop stoke-test > /dev/null
	docker rm stoke-test > /dev/null
}

_recreate_postgres_schema() {
	docker compose exec -it postgres psql -U stoke_user -d stoke -c "drop schema public cascade; create schema public;"
}

_start_client_env() {
	docker compose -f $TASK_DIR/client/client-test-compose.yaml down
	if [[ -n "$ARG_BUILD" ]]
	then
		extra_args="--build"
	fi
	docker compose -f $TASK_DIR/client/client-test-compose.yaml up -d $extra_args
}
