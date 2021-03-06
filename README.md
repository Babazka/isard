isard
=====

Движок для системы мониторинга. Используется **Go** версии 1.11 и выше. 
Для рендера графиков используется [gnuplot](http://www.gnuplotting.org/plotting-data/).

Как пользоваться:
-----------------

* Развернуть ``isard`` локально. Допустим, в директории ``/www/isard``.
* Каждая папка внутри ``/www/isard/data`` целиком описывает какой-то один график. 
* Запустить демона: 
```
make run
```
* По-умолчанию веб-морда доступна по адресу ``http://localhost:7331/``.

Принцип работы:
---------------

* Для каждого графика демон запускается в соответствии с настройками из файлика ``cron.txt``. Эти настройки 
соответствуют синтаксису "часы", "минуты" из ``crontab``
* Для каждого графика выполняется скрипт ``collect.sh``, который вернёт дату и что угодно ещё
* Результат работы скрипта записывается в файл ``data.txt``
* График рендерится по правилам, описанным в ``plot.gnuplot``
* Если графиков очень много и их нужно разбить по категориям, то нужно создать файл ``tags.txt``, в котором 
будут перечислены теги, соответствующие данному графику.
