arguments_test() {
	DESCRIPTION="Run tests"
	SUBCOMMANDS="int|unit"
	INT_DESCRIPTION="Run integration tests"
	INT_OPTIONS="image:i:str"
	UNIT_DESCRIPTION="Run unit tests"
	UNIT_OPTIONS="cover:c:bool html:h:bool func:f:bool name:n:str"
}

task_test() {
	cd $TASK_DIR
	if [[ "$TASK_SUBCOMMAND" == "int" ]]
	then
		if [[ -z "$ARG_IMAGE" ]]
		then
			ARG_IMAGE=stoke-inttest-$(date +%Y%m%d%H%M)
			docker build -t $ARG_IMAGE .
		fi

		echo Using $ARG_IMAGE docker image...

		cd test

		echo Running cert type smoke tests...
		for config in ./configs/cert_types/*sa.yaml
		do
			echo Testing with $config...
			docker run --rm -d --name stoke-cert-test \
				-v $(pwd)/$config:/etc/stoke/config.yaml \
				-v $(pwd)/configs/cert_types/dbinit.yaml:/etc/stoke/dbinit.yaml \
				-p 8080:8080 \
				$ARG_IMAGE -dbinit /etc/stoke/dbinit.yaml

			k6 run k6/smoke/ok_logins.js
			k6 run k6/smoke/bad_logins.js
			docker stop stoke-cert-test
		done

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
