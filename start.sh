#!/bin/sh
#
# 这个文件被用于官方clair镜像的命令入口
#

set -e

cmd="$@"

until psql "host=postgres user=postgres password=password" -c '\q'; do
      >&2 echo "Postgres is unavailable - sleeping"
        sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
