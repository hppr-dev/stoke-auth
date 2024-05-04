# Example Stoke Flask Client Example

How to run:

1. Setup python virtualenv `t venv init -n weapons`
2. Enable virtualenv `t venv enable`
3. Install flask `pip install flask`
4. Install pystokeauth
    a. In development run `pip install -e .` in `STOKE_ROOT/client/pystokeauth`
    b. Otherwise run `pip install pystokeauth`
5. Start a stoke server on localhost:8080
6. Run
    a. With auth enabled    : `flask --app app.py run --port 4000`
    b. Without auth enabled : `STOKE_TEST=yes flask --app app.py run --port 4000`
