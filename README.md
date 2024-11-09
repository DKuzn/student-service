# Student Service

Примерный проект для установки в качестве systemd юнита.

## Документация

### Вывод справки

```bash
dkuzn@localhost:~$ ./student-service --help
An example service to get info about the student.

Usage:
  student [flags]

Flags:
  -n, --fullname string   Full name of the student. (default "Full Name")
  -h, --help              help for student
  -p, --port int          The port is used to listen requests. (default 8080)
```

### Запуск

```bash
dkuzn@localhost:~$ ./student-service --fullname "Ivanov Ivan Ivanovich" --port 8080
```

После запуска единственная конечная точка API доступна по адресу [http://localhost:8080/student](http://localhost:8080/student).

### Маршруты API

#### GET /student

Вывод информацию о студенте, а именно полное имя и PID процесса исполняемого файла. Возвращает JSON сделующего вида.

```json
{
    "proc_id": 12, // ID запущенного процесса
    "full_name": "Ivanov Ivan Ivanovich" // Полное имя студента
}
```

## Порядок настройки в ОС Ubuntu Server

Перейдем в домашнюю директорию пользователя.

```bash
dkuzn@localhost:~$ cd /home/dkuzn
```

Создадим новую директорию и перейдем в нее.

```bash
dkuzn@localhost:~$ mkdir student-service
dkuzn@localhost:~$ cd student-service
```

Получим исполняемый файл из этого репозитория. Для этого необходимо скачать файл по [ссылке](https://github.com/DKuzn/student-service/releases/download/v1/student-service) с помощью иструмента `wget`.

```bash
dkuzn@localhost:~/student-service$ wget https://github.com/DKuzn/student-service/releases/download/v1/student-service
```

По умолчанию, для запуска скачанного исполняемого файла ему необходимо выдать права для запуска. Это сделано для безопасности. Изменить права доступа можно с помощью инстурмента `chmod`.

```bash
dkuzn@localhost:~/student-service$ chmod +x ./student-service
```

Так, команда выше разрешает текущему пользователю запускать исполняемый файл.

Проверим успешность скачивания и выдачи прав доступа тестовым запуском исполняемого файла. В случае успеха выведется сообщение о том, что сервер запущен.

```bash
dkuzn@localhost:~/student-service$ ./student-service
2024/11/09 18:35:01 Student Service is running on port: 8080
```

При создании systemd юнита необходимо описать его конфигурацию. Для этого следует перейти в директорию `/etc/systemd/system`, где находятся файлы с конфигурациями юнитов.

```bash
dkuzn@localhost:~/student-service$ cd /etc/systemd/system
```

Создадим новый файл с конфигурацией нашего systemd юнита. Необходимо обратить внимание на то, что для доступа к операциям с файлами в этой директории необходима привилегия суперпользователя.

```bash
dkuzn@localhost:/etc/systemd/system$ sudo touch student.service
```

Отредактируем этот файл с помощью консольного текстового редактора `vim`. Подробнее с описаниями доступных опций можно ознакомится по [ссылке](https://newadmin.ru/sozdanie-prostogo-systemd-unit/).

```bash
dkuzn@localhost:/etc/systemd/system$ sudo vim student.service
[Unit]
Description=Example systemd service

[Service]
Type=simple
ExecStart=/home/dkuzn/student-service/student-service --fullname "Ivanov Ivan Ivanovich"
Restart=no-failure

[Install]
WantedBy=multi-user.target
~
~
~
:wq
```

Далее активируем наш systemd юнит для автоматического управления системой и запустим его при помощи инстурмента `systemctl`.

```bash
dkuzn@localhost:/etc/systemd/system$ sudo systemctl enable student
dkuzn@localhost:/etc/systemd/system$ sudo systemctl start student
```

С этого момена наш юнит будет автоматически запускаться при старте системы и сервис будет доступен для обработки запросов. Проверим работоспособность сервиса, с помощью инстурмента `curl`. Должна быть выведена строка, похожая на ту, что описана в начале этого файла.

```bash
dkuzn@localhost:/etc/systemd/system$ curl "http://localhost:8080/student"
{"proc_id":814,"full_name":"Ivanov Ivan Ivanovich"}
```

Перезапустим систему, чтобы проверить работоспособность юнита заново.

```bash
dkuzn@localhost:/etc/systemd/system$ sudo reboot
dkuzn@localhost:/etc/systemd/system$ curl "http://localhost:8080/student"
{"proc_id":816,"full_name":"Ivanov Ivan Ivanovich"}
```

После перезапуска системы параметр `proc_id` должен измениться. Это свидетельствует о том, что сервис запустился правильно и его настройка завершена.
