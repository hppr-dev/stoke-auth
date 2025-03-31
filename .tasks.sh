arguments_test() {
	DESCRIPTION="Run tests"
	SUBCOMMANDS="int|unit|client|clean"
	INT_DESCRIPTION="Run integration tests"
	INT_OPTIONS="image:i:str rimage:R:str build:b:bool buildclient:B:bool case:C:str log:l:bool logfile:L:bool progress:P:bool all:a:bool cert:c:bool db:d:bool race:r:bool provider:p:bool env:e:bool"
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
		export STOKE_ADDRESS=localhost:8080
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
			_run_all_configs "race" data_race data_race.yaml data_race.js $ARG_RIMAGE
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

arguments_certs() {
	DESCRIPTION="Manage certificates required for running the example clients"
	SUBCOMMANDS="clean|gen|verify"
}

task_certs() {
	cd $TASK_DIR/client/examples/certs

	if [[ "$TASK_SUBCOMMAND" == "clean" ]]
	then
		echo Removing all client example certificates...
		rm ./*.crt ./*.key
	elif [[ "$TASK_SUBCOMMAND" == "gen" ]]
	then
		for conf in config/*.conf
		do
			cert_name=$(basename ${conf/.conf}).crt
			if [[ ! -f "$cert_name" ]]
			then
				echo Generating $cert_name...
				if [[ "$cert_name" == "ca.crt" ]]
				then
					openssl req -new -x509 -newkey rsa:2048 -config $conf -out $cert_name -days 3650 -extensions v3_req
				else
					openssl req -new -x509 -newkey rsa:2048 -config $conf -out $cert_name -CA ca.crt -CAkey ca.key -days 3650 -extensions v3_req
				fi
			else
				echo Found $cert_name, skipping...
			fi
		done

	elif [[ "$TASK_SUBCOMMAND" == "verify" ]]
	then
		openssl verify -CAfile ca.crt *.crt 
	fi

}

arguments_clients() {
	DESCRIPTION="Manage test/example client docker containers"
	SUBCOMMANDS="up|down|sh|ps|logs"
	CLIENTS_OPTIONS="build:b:bool detach:d:bool"
	SH_REQUIREMENTS="service:s:str"
}

task_clients() {
	if [[ "$TASK_SUBCOMMAND" == "up" ]] && ! docker ps | grep stoke_server > /dev/null
	then
		echo Could not find running stoke_server. Please run task stoke up first
		exit 1
	fi
	_compose_task "$TASK_DIR/client/client-test-compose.yaml"
}

arguments_stoke() {
	DESCRIPTION="Manage stoke compose containers"
	SUBCOMMANDS="up|down|sh|ps|logs"
	STOKE_OPTIONS="build:b:bool detach:d:bool"
	SH_REQUIREMENTS="service:s:str"
}

task_stoke() {
	_compose_task "$TASK_DIR/client/examples/stoke-server/docker-compose.yaml"
}

arguments_kube() {
	DESCRIPTION="Kubernetes helper scripts"
	SUBCOMMANDS="pvc|tkn"
	PVC_OPTIONS="create:c:bool update:u:bool delete:d:bool"
	TKN_OPTIONS="local:l:bool git:g:bool"
}

task_kube() {
	if [[ "$TASK_SUBCOMMAND" == "pvc" ]]
	then
		if [[ -n "$ARG_CREATE" ]]
		then
			cat << EOF| kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: stoke-local
spec:
  storageClassName: local-path
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
EOF

		fi

		if [[ -n "$ARG_UPDATE" ]]
		then
			echo "Starting transfer pod..."
			cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: stoke-xfer-pod
spec:
  restartPolicy: Never
  volumes:
    - name: stoke
      persistentVolumeClaim:
        claimName: stoke-local
  containers:
    - name: xfer
      image: alpine:latest
      command: ["tail", "-f", "/dev/null"]
      volumeMounts:
        - name: stoke
          mountPath: /mnt/stoke
EOF

			while [[ "$(kubectl get pod stoke-xfer-pod -o 'jsonpath={.status.phase}')" != "Running" ]]
			do
				echo "Waiting for pod to be ready..."
				sleep 1
			done

			echo "Deleting existing files..."
			kubectl exec -i stoke-xfer-pod -- /bin/ash -c 'rm -rf /mnt/stoke/*'

			echo "Updating to latest files..."
			cd $TASK_DIR
			tar -cf - . | kubectl exec -i stoke-xfer-pod -- tar xf - -C /mnt/stoke

			echo "List of files in pvc:"
			kubectl exec -i stoke-xfer-pod -- ls /mnt/stoke

			echo "Removing xfer pod..."
			kubectl delete pod stoke-xfer-pod --now

		fi
		if [[ -n "$ARG_DELETE" ]]
		then
			echo "Deleting PVC..."
			kubectl delete pvc stoke-local
		fi

	fi
	if [[ "$TASK_SUBCOMMAND" == "tkn" ]]
	then
		if [[ -n "$ARG_LOCAL" ]]
		then
			echo Running pipeline with local pvc...
			kubectl create -f $TASK_DIR/test/tekton/local_run.yaml
		fi
		if [[ -n "$ARG_GIT" ]]
		then
			echo Running pipeline on latest from git...
			kubectl create -f $TASK_DIR/test/tekton/git_run.yaml
		fi
	fi
	echo "Done."

}

_compose_task() { #compose_file
	compose_file=$1
	if [[ "$TASK_SUBCOMMAND" == "up" ]]
	then
		if [[ -n "$ARG_BUILD" ]]
		then
			extra="--build"
		fi
		if [[ -n "$ARG_DETACH" ]]
		then
			extra="$extra -d"
		fi
		docker compose -f $compose_file up $extra

	elif [[ "$TASK_SUBCOMMAND" == "down" ]]
	then
		docker compose -f $compose_file down

	elif [[ "$TASK_SUBCOMMAND" == "ps" ]]
	then
		docker compose -f $compose_file ps

	elif [[ "$TASK_SUBCOMMAND" == "sh" ]]
	then
		docker compose -f $compose_file exec -it $ARG_SERVICE /bin/sh

	elif [[ "$TASK_SUBCOMMAND" == "logs" ]]
	then
		docker compose -f $compose_file logs -f
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
		-v $TASK_DIR/client/examples/certs/stoke.crt:/etc/stoke/stoke.crt \
		-v $TASK_DIR/client/examples/certs/stoke.key:/etc/stoke/stoke.key \
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
	if ! k6 run --insecure-skip-tls-verify $k6_args k6/$k6file
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
	if [[ -n "$ARG_BUILDCLIENT" ]]
	then
		extra_args="--build"
	fi
	docker compose -f $TASK_DIR/client/client-test-compose.yaml up -d $extra_args
}
