Используя средства языка Scheme для метапрограммирования, реализуйте каркас поддержки типа данных «структура» («запись»). Пусть объявление нового типа «структура» осуществляется с помощью вызова (define-struct тип-структуры (имя-поля-1 имя-поля-2 ... имя-поля-n). Тогда после объявления структуры программисту становятся доступны:

    Процедура — конструктор структуры вида (make-тип-структуры значение-поля-1 значение-поля-2 ... значение-поля-n), возвращающий структуру, поля которой инициализированы перечисленными значениями.
    Предикат типа вида (тип-структуры? объект), возврщающая #t если объект является структурой типа тип-структуры и #f в противном случае.
    Процедуры для получения значения каждого из полей структуры вида (тип-структуры-имя-поля объект).
    Процедуры модификации каждого из полей структуры вида (set-тип-структуры-имя-поля! объект новое-значение). 

Пример использования каркаса:

```LISP
(define-struct pos (row col)) ; Объявление типа pos
(define p (make-pos 1 2))     ; Создание значения типа pos

(pos? p)    ⇒ #t

(pos-row p) ⇒ 1
(pos-col p) ⇒ 2

(set-pos-row! p 3) ; Изменение значения в поле row
(set-pos-col! p 4) ; Изменение значения в поле col

(pos-row p) ⇒ 3
(pos-col p) ⇒ 4
```