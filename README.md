# multi-ssh

读取 ansible Inventory 文件，批量执行命令。

# 编译

```bash
$ go env -w GOOS=linux
$ go build -o multi-ssh
```

# 使用方法

```bash
$ ./multi-ssh -h
Usage of ./multi-ssh:
  -H string
        Ansible Inventory File (default "hosts")
  -cmd string
        Command To Execute
  -debug
        Show Debug Info
  -g string
        Hosts In Which Group
  -i string
        PrimaryKey File (default "/home/yaoliang/.ssh/id_rsa")
  -try
        Try Run, Just Show Command
```

## 参数解析

`-H` 指定 ansible inventory 文件路径，默认搜索路径如下：

- 当前目录下的`hosts`文件
- `~/ansible/hosts`文件
- `/etc/ansible/hosts`文件

`-g` 指定要执行命令的主机组

`-cmd` 指定要执行的命令，可能需要使用`""`包裹

`-i` 指定 ssh 私钥路径，默认搜索路径如下：

- `~/.ssh/id_rsa`文件

`-try` 尝试执行，只打印命令，不执行
`-debug` 显示调试信息

## 示例

```bash
$ ./multi-ssh -g www -cmd "df -Th /dev/shm"
2025/03/18 11:23:04 [SUCCESS]   www     [0/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs   23G   21G  2.3G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [9/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  141G  127G   15G  90% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [8/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  141G  127G   15G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [2/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs   94G   85G  9.5G  90% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [10/11] 10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs   94G   85G  9.4G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [4/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  141G  127G   15G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [5/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs   94G   85G  9.4G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [7/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  140G  127G   14G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [1/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs   94G   85G  9.4G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [3/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  141G  127G   15G  91% /dev/shm

2025/03/18 11:23:04 [SUCCESS]   www     [6/11]  10.x.x.x
Filesystem     Type   Size  Used Avail Use% Mounted on
tmpfs          tmpfs  140G  127G   14G  91% /dev/shm
```
