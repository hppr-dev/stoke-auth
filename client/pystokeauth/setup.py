# Always prefer setuptools over distutils
from setuptools import setup, find_packages
import pathlib

here = pathlib.Path(__file__).parent.resolve()

long_description = (here / "README.md").read_text(encoding="utf-8")

setup(
    name="pystokeauth",
    version="0.1.0",
    description="Client library to use JWTs issued by a stoke server",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/hppr-dev/stoke-auth",
    author="Stephen Walker",  # Optional
    author_email="swalker@hppr.dev",  # Optional
    classifiers=[  # Optional
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3 :: Only",
    ],
    keywords="authentication, jwt, microservice",
    #package_dir={"": "src"},
    packages=find_packages(where="."),
    python_requires=">=3.7, <4",
    install_requires=["pyjwt[crypto]"],
    extras_require={
        "flask": ["flask"],
        "django": ["django"],
    },
    project_urls={  # Optional
        "Bug Reports": "https://github.com/hppr-dev/stoke-auth/issues",
        "Source": "https://github.com/hppr-dev/stoke-auth/",
    },
)
