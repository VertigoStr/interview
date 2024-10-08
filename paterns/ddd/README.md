# Предметно-ориентированное проектирование

## Агрегаты
Агрегаты представляют из себя автономную единицу, обладающую состояниями. Агрегат это дерево сущностей 
(может включать в себя другие сущности).

Например, для интернет-магазина, "заказ" - хороший пример агрегата. Жизненный цикл агрегатов связан с изменением их состояний.

### Свойства агрегатов

Агрегаты подразумевают использование транзакционных систем в реляционных базах данных.

Инварианты - неизменяемые свойства, которые должны выполняться при изменении системы. 
(Могут быть нарушены при определенной настройке транзакций)

Аномалии - (тема транзакций) аномалия, которая возникает (один из примеров) при ситуации когда мы из двух разных транзакций пытаемся изменить одни и те же данные 
(например, ячейки А и В) так, что первая транзакция затронет ячейку A, а вторая — ячейку В и мы получим неконсистентный результат.

Оптимистичная блокировка - (тема транзакций) настройка уровня изоляций, позволяющий избавиться от возможных аномалий.

### Как хранятся агрегаты в базе данных
На примере `mongoDB` данные хранятся как документ с деревом сущностей.
Для изменения подсущности - достается все дерево.
![img](img.png)

Чем больше сущностей мы объединим в один агрегат, тем больше мы потребляем памяти и ресурсов.


### Transactional outbox pattern
Типы сообщений:
1. `Comands` - операции, обновляющие наши агрегаты (`create`, `update`). Включает в себя *валидацию*, *инварианты*.
2. `Query` - операции на чтение данных (`get`, `find`, `filter`).
3. `Events` - операции, описывающие факт совершенного изменения (`has been created`).

`Comands` могут порождать `Events` (Одна команда -> один ивент ; Одна команда -> много ивентов).

Проблема:
Есть command над агрегатом, он порождает update в базу данных и event для мессадж брокера.
Ивенты, попадающие в мессадж брокер могут быть не в том порядке или вообще не прийти (если упал сервис, моргнула сеть).

Решение:
Необходимо свойство атомарности между обновлениями. 

`Transactional outbox pattern` подразумевает в рамках 1 транзакции создать event для DB и создать event для message broker.
![img_1.png](img_1.png)

Как это выглядит для базы данных:
Нужна таблица агрегатов и таблица ивентов, а так же база данных должа поддерживать ACID:
![img_2.png](img_2.png)

Можно использовать атомарную запись документов и убрать из базы данных таблицу агрегатов 
(тогда получим Event sourcing систему):
![img_3.png](img_3.png)


## Event Store
![img_4.png](img_4.png)
База данных ивентов, обладает функционалом:
- Запросить ивенты для конкретного объекта/агрегата (по id)
- Поместить в базу данных ивент (обновление агрегата)
- Не обязан поддерживать удаление/изменение событий.

`EventLog` - коллекция ивентов для конкретного агрегата (с определенным id)
`EventRecords` - мы храним ивент обернутый в произвольную структуру,
со своей метаинформацией.

`EventRecords` хранит в себе:
- `Метаинформацию` (Для десереализации можно хранить имя ивента или имя класса, в который парсится ивент, 
также там может хранится: номер ивента в последовательности, timestamp, user_id, version и т.д.)
- `Event`: хранится как атрибут в реляционной таблице в виде атомарных струткур: строка/массив байтов/сериализованный протобуф и т.д.

Номер ивента в последовательности:
Важно чтобы события получались по порядку. 
В реляционных бд - юзаем Serial. В mongoDB - пишем кастомную функцию, которая генерит id в зависимости от принадлежности текущего ивента.
