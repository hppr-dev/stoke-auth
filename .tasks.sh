arguments_test() {
	DESCRIPTION="Run tests"
	SUBCOMMANDS="int|unit|client|clean"
	INT_DESCRIPTION="Run integration tests"
	INT_OPTIONS="image:i:str build:b:str case:C:str log:l:bool progress:P:bool all:a:bool cert:c:bool db:d:bool race:r:bool provider:p:bool env:e:bool"
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
			echo Running with most found image: $ARG_IMAGE
		fi

		if [[ -n "$ARG_BUILD" ]] || [[ -z "$ARG_IMAGE" ]]
		then
			ARG_IMAGE=stoke-inttest-$(date +%Y%m%d%H%M)
			echo Building new image $ARG_IMAGE...
			docker build -t $ARG_IMAGE .
		fi

		echo Starting supplemental containers...
		cd $TASK_DIR/compose
		docker compose --profile integration up -d

		cd $TASK_DIR/test

		if [[ -z "$ARG_CERT$ARG_DB$ARG_RACE$ARG_PROVIDER$ARG_ENV" ]]
		then
			ARG_ALL="yes"
		fi

		if [[ -n "$ARG_CERT$ARG_ALL" ]]
		then
			echo ====================================================== Running cert smoke tests...
			for config in ./configs/cert_types/*
			do
				_run_k6_test $config smoke_test.yaml smoke.js
			done
		fi


		if [[ -n "$ARG_DB$ARG_ALL" ]]
		then
			echo ====================================================== Running database smoke tests...
			for config in ./configs/database_types/*
			do
				_run_k6_test $config smoke_test.yaml smoke.js
			done
		fi

		if [[ -n "$ARG_RACE$ARG_ALL" ]]
		then
			echo Running data race tests...
			echo TODO TODO TODO
		fi

		if [[ -n "$ARG_PROVIDER$ARG_ALL" ]]
		then
			echo Running provider tests...
			echo TODO TODO TODO
		fi

		if [[ -n "$ARG_ENV$ARG_ALL" ]]
		then
			echo Running client environment tests...
			echo TODO TODO TODO
		fi

		cd $TASK_DIR/compose
		docker compose --profile=integration down

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
		cd compose
		docker compose --profile integration down 

		docker stop stoke-test
		docker rm stoke-test

		if [[ -f "$TASK_DIR/cover.html" ]]
		then
			rm "$TASK_DIR/cover.html"
		fi

		if [[ -f "$TASK_DIR/cover.out" ]]
		then
			rm "$TASK_DIR/cover.out"
		fi
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


_run_k6_test() { # config dbinit k6file
	# Config file relative to CWD
	config=$1
	# dbinit file relative to CWD/configs/dbinit/
	dbinit=$2
	# k6 file relative to CWD/k6/
	k6file=$3

	if [[ -n "$ARG_CASE" ]] && ! [[ $config =~ "$ARG_CASE" ]]
	then
		echo Did not match case, skipping $config...
		continue
	fi

	docker_args="--name stoke-test \
		-v $(pwd)/$config:/etc/stoke/config.yaml \
		-v $(pwd)/configs/dbinit/$dbinit:/etc/stoke/dbinit.yaml \
		-p 8080:8080 \
		$ARG_IMAGE -dbinit /etc/stoke/dbinit.yaml"
	docker run --rm -d $docker_args > /dev/null
	sleep 1

	if ! docker ps | grep stoke-test > /dev/null
	then
		echo Could not start container with $config. Running again to show output.
		docker run $docker_args
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
	k6 run $k6_args k6/$k6file
	docker stop stoke-test > /dev/null
}
