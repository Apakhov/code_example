(use-syntax (ice-9 syncase))

; <цифра> ::= 0|1|2|3|4|5|6|7|8|9
; <знак> ::= -|+|*|/|^
; <скобка> ::= ( | )
; <литер> ::= a|b|c|...|y|z|A|B|C|...|Y|Z
; <переменная> ::=  <литер> | <литер><переменная>
; <об.число> ::= <цифра> | <цифра><об.число>
; <дес.число> ::= <об.число>.<об.число>
; <число> ::= <об.число> | <дес.число>
; <эксп.число> ::= <число>e<число> | <число>E<число> | <число>e+<число> | <число>E+<число> | <число>e-<число> | <число>E-<число>

(define-syntax make-source
  (syntax-rules ()
    ((_ smth) (cond
                ((vector? smth) (vector 'sequence #f smth 0 (vector-length smth)))
                ((string? smth) (vector 'sequence #f (list->vector (string->list smth)) 0 (vector-length (list->vector (string->list smth)))))
                ((list? smth) (vector 'sequence #f (list->vector smth) 0 (vector-length (list->vector smth))))
                (else 'not-a-sequence)))
    ((_ smth end) (cond
                    ((vector? smth) (vector 'sequence end smth 0 (vector-length smth)))
                    ((string? smth) (vector 'sequence end (list->vector (string->list smth)) 0 (vector-length (list->vector (string->list smth)))))
                    ((list? smth) (vector 'sequence end (list->vector smth) 0 (vector-length (list->vector smth))))
                    (else 'not-a-sequence)))))

(define (peek sequence)
  (if (= (vector-ref sequence 3) (vector-ref sequence 4))
      (vector-ref sequence 1)
      (vector-ref (vector-ref sequence 2) (vector-ref sequence 3))))

(define (next sequence)
  (if (= (vector-ref sequence 3) (vector-ref sequence 4))
      (vector-ref sequence 1)
      (begin
        (vector-set! sequence 3 (+ (vector-ref sequence 3) 1))
        (vector-ref (vector-ref sequence 2) (- (vector-ref sequence 3) 1)))))

(define (tokenize str)
  (define source (make-source str 'eof))
  (define (checker xs ys)
    (if (null? xs)
        (reverse ys)
        (and (car xs)
             (checker (cdr xs) (cons (car xs) ys)))))
  (define (get-num source lis)
    (if (and (or (equal? (peek source) #\+)
                 (equal? (peek source) #\-)) (not (null? lis)) (or (equal? #\e (car lis)) (equal? #\E (car lis))))
        (get-num source (cons (next source) lis))
        (if (and (not (equal? (peek source) 'eof)) (not (char-whitespace? (peek source))) (not (or (equal? (peek source) #\^)
                                                                                                   (equal? (peek source) #\*)
                                                                                                   (equal? (peek source) #\/)
                                                                                                   (equal? (peek source) #\+)
                                                                                                   (equal? (peek source) #\-)
                                                                                                   (equal? (peek source) #\)))))
            (get-num source (cons (next source) lis))          
            (string->number (list->string (reverse lis))))))
  (define (get-symbol source sym)
    (if (and (not (equal? (peek source) 'eof)) (or (and (>= (char->integer (peek source)) (char->integer #\a))
                                                        (<= (char->integer (peek source)) (char->integer #\z)))
                                                   (and (>= (char->integer (peek source)) (char->integer #\A))
                                                        (<= (char->integer (peek source)) (char->integer #\Z)))))
        (get-symbol source (cons (next source) sym))
        (string->symbol (list->string (reverse sym)))))
  (define (iter source outs)
    (cond
      ((equal? (peek source) 'eof) (reverse outs))
      ((and (>= (char->integer (peek source)) (char->integer #\0))
            (<= (char->integer (peek source)) (char->integer #\9)))
       (iter source (cons (get-num source '()) outs)))
      ((or (and (>= (char->integer (peek source)) (char->integer #\a))
                (<= (char->integer (peek source)) (char->integer #\z)))
           (and (>= (char->integer (peek source)) (char->integer #\A))
                (<= (char->integer (peek source)) (char->integer #\Z))))
       (iter source (cons (get-symbol source '()) outs)))
      ((char-whitespace? (peek source)) (begin
                                          (next source)
                                          (iter source outs)))
      ((or (equal? (peek source) #\^)
           (equal? (peek source) #\*)
           (equal? (peek source) #\/)
           (equal? (peek source) #\+)
           (equal? (peek source) #\-))       
       (iter source (cons (string->symbol (list->string (list (next source)))) outs)))
      ((or (equal? (peek source) #\()
           (equal? (peek source) #\)))
       (iter source (cons (list->string (list (next source))) outs)))
      (else '(#f))))
  (checker (iter source '()) '()))

(define (parse xs)
  (define (correct? xs)
    (define (oper? x)
      (or (equal? x '^)
          (equal? x '*)
          (equal? x '/)
          (equal? x '+)
          (equal? x '-)))
    (define (-ORsymb? x)
      (or (equal? '- x) (not (oper? x))))
    (define (symb? x)
      (and (not (oper? x)) (not (equal? x ")"))))
    (define (<>amount-checker xs count)
      (cond
        ((null? xs) (= count 0))
        ((equal? "(" (car xs)) (<>amount-checker (cdr xs) (+ count 1)))
        ((equal? ")" (car xs)) (<>amount-checker (cdr xs) (- count 1)))
        (else (<>amount-checker (cdr xs) count))))
    (define (position-checker xs right-term?)
      (cond
        ((null? xs) #t)
        ((equal? (car xs) "(")  (position-checker (cdr xs) -ORsymb?))
        ((and (equal? (car xs) ")") (equal? right-term? oper?)) (position-checker (cdr xs) oper?))
        ((and (oper? (car xs)) (right-term? (car xs)))  (position-checker (cdr xs) symb?))
        ((and (symb? (car xs)) (right-term? (car xs)))  (position-checker (cdr xs) oper?))
        (else #f)))
    (and (<>amount-checker xs 0) (position-checker xs -ORsymb?)))

  (define (<>-destroyer xs)
    (define (iter left-part right-part)
      (cond
        ((null? right-part) (reverse left-part))
        ((null? left-part) (iter `(,(car right-part)) (cdr right-part)))
        ((and (equal? (car left-part) "(") (equal? (cadr right-part) ")")) (iter (cdr left-part) (cons (car right-part) (cddr right-part))))
        (else (iter (cons (car right-part) left-part) (cdr right-part)))))
    (and xs (iter '() xs)))

  (define (less^ xs)
    (define (iter left-part right-part cur-depth res res-depth)
      (cond
        ((null? right-part) (cons res res-depth))
        ((equal? "(" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (+ cur-depth 1) res res-depth))
        ((equal? ")" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (- cur-depth 1) res res-depth))
        ((and (equal? '^ (car right-part)) (>= cur-depth res-depth)) (iter (cons (cadr right-part) (cons (car right-part) left-part)) (cddr right-part) cur-depth (append (reverse (cdr left-part)) `((,(car left-part) ^ ,(cadr right-part)))  (cddr right-part)) cur-depth))                                             
        (else (iter (cons (car right-part) left-part) (cdr right-part) cur-depth res res-depth))))
    (iter '() xs 0 #f -1))

  (define (less+- xs)
    (define (iter left-part right-part cur-depth res res-depth)
      (cond
        ((null? right-part) (cons (and res (reverse res)) res-depth))
        ((equal? ")" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (+ cur-depth 1) res res-depth))
        ((equal? "(" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (- cur-depth 1) res res-depth))
        ((and (or (equal? '+ (car right-part)) (equal? '- (car right-part))) (>= (length right-part) 2) (>= cur-depth res-depth))       
         (iter (cons (cadr right-part) (cons (car right-part) left-part)) (cddr right-part) cur-depth (append (reverse (cdr left-part)) `((,(cadr right-part) ,(car right-part) ,(car left-part)))  (cddr right-part)) cur-depth))                                            
        (else (iter (cons (car right-part) left-part) (cdr right-part) cur-depth res res-depth))))
    (iter '() (reverse xs) 0 #f -1))

  (define (less*/ xs)
    (define (iter left-part right-part cur-depth res res-depth)
      (cond
        ((null? right-part) (cons (and res (reverse res)) res-depth))
        ((equal? ")" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (+ cur-depth 1) res res-depth))
        ((equal? "(" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (- cur-depth 1) res res-depth))
        ((and (or (equal? '* (car right-part)) (equal? '/ (car right-part))) (>= cur-depth res-depth)) (iter (cons (cadr right-part) (cons (car right-part) left-part)) (cddr right-part) cur-depth (append (reverse (cdr left-part)) `((,(cadr right-part) ,(car right-part) ,(car left-part)))  (cddr right-part)) cur-depth))                                             
        (else (iter (cons (car right-part) left-part) (cdr right-part) cur-depth res res-depth))))
    (iter '() (reverse xs) 0 #f -1))

  (define (less-unary- xs)
    (define (iter left-part right-part cur-depth)
      (cond
        ((null? right-part) '(#f . -1))
        ((equal? "(" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (+ cur-depth 1)))
        ((equal? ")" (car right-part)) (iter (cons (car right-part) left-part) (cdr right-part) (- cur-depth 1)))
        ((and (equal? '- (car right-part)) (or (null? left-part) (equal? "(" (car left-part))) (not (equal? "(" (cadr right-part))))
         `(,(append (reverse left-part) `((- ,(cadr right-part))) (cddr right-part)) . ,cur-depth))                                             
        (else (iter (cons (car right-part) left-part) (cdr right-part) cur-depth))))
    (iter '() xs 0))

  (define (main xs)
    (define (max xs res res-depth)    
      (if (null? xs)
          res
          (if (> (cdr (car xs)) res-depth)
              (max (cdr xs) (caar xs) (cdar xs))
              (max (cdr xs) res res-depth))))
    (and xs (correct? xs) (not (null? xs)) (if (null? (cdr xs))
                                               (car xs)
                                               (and (max `(,(less^ xs) ,(less*/ xs) ,(less-unary- xs) ,(less+- xs))  #f -1)
                                                    (main (<>-destroyer (max `(,(less^ xs) ,(less*/ xs) ,(less-unary- xs) ,(less+- xs)) #f -1)))))))
  (main xs))



(define (tree->scheme xs)
  (and xs (if (list? xs)
              (if (= (length xs) 3)
                  `(,(tree->scheme (cadr xs)) ,(tree->scheme (car xs)) ,(tree->scheme (caddr xs)))
                  `(,(tree->scheme (car xs)) ,(tree->scheme (cadr xs))))
              (if (equal? xs '^)
                  'expt
                  xs))))



