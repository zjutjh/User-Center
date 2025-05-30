# User-Center-grpc
精弘用户中心

### Build

Before building the project, you need to install the clang-format. And then run the following command to install the
grpc tool and build the project.

### Install grpc binary

```bash
make install
```

### Generate proto files

```bash
make genproto
```

### Build Apiserver

```bash
make apiserver
```

## Run

```bash
make run
```