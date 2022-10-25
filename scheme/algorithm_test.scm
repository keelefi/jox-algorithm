#!/usr/bin/guile -s
!#

(use-modules (srfi srfi-11)
             (srfi srfi-64))

(use-modules (json))

(define (jobs-alist-equal? jobs-a jobs-b)
    (define (sort-alist alist)
        (sort alist (lambda (a b) (string< (car a) (car b)))))

    (define (equal-jobs? job-a job-b)
        (let ((job-a-name (car job-a))
              (job-b-name (car job-b))
              (job-a-after (sort (cdr (assoc "after" (cdr job-a))) string<))
              (job-b-after (sort (cdr (assoc "after" (cdr job-b))) string<))
              (job-a-before (sort (cdr (assoc "before" (cdr job-a))) string<))
              (job-b-before (sort (cdr (assoc "before" (cdr job-b))) string<)))
            (and
                (equal? job-a-name job-b-name)
                (equal? job-a-after job-b-after)
                (equal? job-a-before job-b-before))))

    (define (jobs-alist-sorted-equal? jobs-a-sorted jobs-b-sorted)
        (cond
            ((null? jobs-a-sorted) (null? jobs-b-sorted))
            ((null? jobs-b-sorted) #f)
            (else
                (let ((current-a (car jobs-a-sorted))
                      (current-b (car jobs-b-sorted)))
                    (if (not (equal-jobs? current-a current-b))
                        #f
                        (jobs-alist-sorted-equal? (cdr jobs-a-sorted) (cdr jobs-b-sorted)))))))

    (if (not (= (length jobs-a) (length jobs-b)))
        #f
        (let ((jobs-a-sorted (sort-alist jobs-a))
              (jobs-b-sorted (sort-alist jobs-b)))
            (jobs-alist-sorted-equal? jobs-a-sorted jobs-b-sorted))))

(include "algorithm.scm")

(define (get-test-file-names)
    (define result '())
    (define test-file-dir "../tests")
    (define dir-obj (opendir test-file-dir))
    (do ((file-entry (readdir dir-obj) (readdir dir-obj)))
        ((eof-object? file-entry))
        (set! result (cons file-entry result)))
    (closedir dir-obj)
    result)

(define (is-valid-file basename)
    (let ((filename (string-append "../tests/" basename)))
        (cond
            ((not (file-exists? filename)) #f)
            ((equal? "." (string-take basename 1)) #f)
            ((not (equal? ".json" (string-take-right basename 5))) #f)
            (else #t))))

(define (is-error error-obj)
    (equal? "ERROR" (cdr (assoc "type" error-obj))))
(define (is-warning error-obj)
    (equal? "WARNING" (cdr (assoc "type" error-obj))))

(define (recreate-jobs jobs)
    (define (recreate-job job)
        (let* ((job-name (car job))
              (job-value (cdr job))
              (after-array (assoc "after" job-value))
              (before-array (assoc "before" job-value))
              (after-list (if after-array (array->list (cdr after-array)) '()))
              (before-list (if before-array (array->list (cdr before-array)) '())))
            `(,job-name . (("after" . ,after-list) ("before" . ,before-list)))))
    (let ((key (car jobs))
          (alist (cdr jobs)))
        (cons key (map recreate-job alist))))

(define (print-job job)
    (let* ((job-name (car job))
           (job-value (cdr job))
           (after (cdr (assoc "after" job-value)))
           (before (cdr (assoc "before" job-value))))
        (display "name: ") (display job-name) (newline)
        (display "after: ") (display after) (newline)
        (display "before: ") (display before) (newline)))

(define (run-test-case filename)
    (define json-document
        (with-input-from-file (string-append "../tests/" filename)
            (lambda () (json->scm (current-input-port)))))

    (define jobs-input (recreate-jobs (assoc "input" json-document)))
    (define jobs-output (recreate-jobs (assoc "output" json-document)))

    (define targets (array->list (cdr (assoc "targets" json-document))))
    (define errors (array->list (cdr (assoc "errors" json-document))))

    (define output-expected (cdr jobs-output))
    (define warnings-expected (filter is-warning errors))
    (define error-expected (filter is-error errors))

    (cond 
        ((null? error-expected)
            (let-values (((output-actual warnings-actual) (algorithm (cdr jobs-input) targets)))
                (test-assert "output" (jobs-alist-equal? output-expected output-actual))
                (test-equal "warnings" warnings-expected warnings-actual)))
        (else
            (with-exception-handler
                (lambda (exception-obj)
                    (test-equal "exception message"
                        (cdr (assoc "message" (car error-expected)))
                        (algorithm-exception-message exception-obj)))
                (lambda ()
                    (algorithm (cdr jobs-input) targets)
                    (test-assert "no exception was thrown" #f))
                #:unwind? #t
                #:unwind-for-type &algorithm-exception))))

(define (make-test-case-name filename)
    (string-append "test-" (basename filename ".json")))

(define test-filenames (filter is-valid-file (get-test-file-names)))

(define (filename-to-test-result filename)
    (let ((test-case-name (make-test-case-name filename)))
        (test-runner-current (test-runner-create))
        (test-begin test-case-name)
        (run-test-case filename)
        (test-end test-case-name)
        (test-runner-fail-count (test-runner-get))))

(define test-results (map filename-to-test-result test-filenames))

;;; fail-count can never be negative. sum the results together and see if they're greater than zero.
(if (< 0 (apply + test-results))
    (exit 1)
    (exit 0))
