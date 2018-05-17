(use-syntax (ice-9 syncase))

(define-syntax define-struct
  (syntax-rules ()
    ((_ struct-name (x . xs))
     (begin
       (eval
        `(define (,(string->symbol (string-append "make-" (symbol->string 'struct-name))) x . xs)
           (cons ',(string->symbol (string-append "structure-" (symbol->string 'struct-name))) (map cons '(x . xs) (list x . xs))))
        (interaction-environment))
       (eval
        `(define (,(string->symbol (string-append (symbol->string 'struct-name) "?")) smth)
           (and (list? smth) (> (length smth) 1) (equal? (car smth) ',(string->symbol (string-append "structure-" (symbol->string 'struct-name))))))
        (interaction-environment))
       (eval
        (letrec ((constructor (lambda (fields res)
                                (if (null? fields)
                                    res
                                    (constructor
                                     (cdr fields)
                                     (append res `((define (,(string->symbol (string-append "set-" (symbol->string 'struct-name) "-" (symbol->string (car fields)) "!")) struct new)
                                                     (define (replacer struct field new)
                                                       (if (equal? field (car (car struct)))
                                                           (cons (cons field new) (cdr struct))
                                                           (cons (car struct) (replacer (cdr struct) field new))))
                                                     (if (and (list? struct) (not (null? struct)) (equal? (car struct)
                                                                                                          (string->symbol (string-append "structure-" (symbol->string 'struct-name)))))                                  
                                                         (set-cdr! struct (replacer (cdr struct) ',(car fields) new))
                                                         '(ERROR: not struct-name))))))))))
          (constructor '(x . xs) '(begin)))
        (interaction-environment))
       (eval
        (letrec ((constructor (lambda (fields res)
                                (if (null? fields)
                                    res
                                    (constructor
                                     (cdr fields)
                                     (append res `((define (,(string->symbol (string-append (symbol->string 'struct-name) "-" (symbol->string (car fields)))) struct)
                                                     (if (and (list? struct) (not (null? struct)) (equal? (car struct)
                                                                                                          (string->symbol (string-append "structure-" (symbol->string 'struct-name)))))                                  
                                                         (cdr (assoc ',(car fields) (cdr struct)))
                                                         '(ERROR: not struct-name))))))))))
          (constructor '(x . xs) '(begin)))
        (interaction-environment))))))