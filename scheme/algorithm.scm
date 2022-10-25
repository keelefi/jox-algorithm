#!/usr/bin/guile -s
!#

(use-modules
    (ice-9 exceptions)
    (ice-9 format)
    (srfi srfi-1)
    (srfi srfi-9))

(define-record-type job
    (make-job name after before)
    job?
    (name job-name)
    (after job-after set-job-after!)
    (before job-before set-job-before!))

(define-exception-type
    &algorithm-exception        ; name
    &error                      ; parent
    make-algorithm-exception    ; constructor
    algorithm-exception?        ; predicate
    (message algorithm-exception-message))

(define-exception-type
    &cyclic-dependency-exception
    &algorithm-exception
    make-cyclic-dependency-exception-internal
    cyclic-dependency-exception?)
(define (make-cyclic-dependency-exception)
    (make-cyclic-dependency-exception-internal "Cyclic dependency detected"))

(define-exception-type
    &job-depends-on-itself-exception
    &algorithm-exception
    make-job-depends-on-itself-exception-internal
    job-depends-on-itself-exception?)
(define (make-job-depends-on-itself-exception job-name)
    (make-job-depends-on-itself-exception-internal (format #f "Job '~a' depends on itself" job-name)))

(define-exception-type
    &no-targets-exception
    &algorithm-exception
    make-no-targets-exception-internal
    no-targets-exception?)
(define (make-no-targets-exception)
    (make-no-targets-exception-internal "No targets"))

(define-exception-type
    &job-not-found-exception
    &algorithm-exception
    make-job-not-found-exception-internal
    job-not-found-exception?)
(define (make-job-not-found-exception job-depender job-dependee)
    (make-job-not-found-exception-internal
        (format #f "Job '~a' references job '~a', but job '~a' does not exist" job-depender job-dependee job-dependee)))

(define-exception-type
    &target-not-found-exception
    &algorithm-exception
    make-target-not-found-exception-internal
    target-not-found-exception?)
(define (make-target-not-found-exception job-name)
    (make-target-not-found-exception-internal
        (format #f "'targets' references job '~a', but job '~a' does not exist" job-name job-name)))

(define (make-warning-job-not-needed job-name)
    `(("message" . ,(format #f "Job '~a' is not required" job-name)) ("type" . "WARNING") ("error" . "JOB NOT REQUIRED")))

(define (algorithm jobs-input targets)
    ; set all after-before relations
    (define (set-after-before-pairs jobs-input job)
        (define (get-job-name job-input)
            (car job-input))
        (define (job-after job-name job-input)
            (member job-name (cdr (assoc "after" (cdr job-input)))))
        (define (job-before job-name job-input)
            (member job-name (cdr (assoc "before" (cdr job-input)))))
        (let* ((job-name (car job))
               (after-list-original (cdr (assoc "after" (cdr job))))
               (after-list-add (map get-job-name (filter (lambda (job-input) (job-before job-name job-input)) jobs-input)))
               (after-list (append after-list-original after-list-add))
               (before-list-original (cdr (assoc "before" (cdr job))))
               (before-list-add (map get-job-name (filter (lambda (job-input) (job-after job-name job-input)) jobs-input)))
               (before-list (append before-list-original before-list-add)))
            `(,job-name . (("after" . ,after-list) ("before" . ,before-list)))))
    (define jobs (map (lambda (job) (set-after-before-pairs jobs-input job)) jobs-input))

    ; check if a job depends on itself
    (define (job-depends-on-itself job)
        (let ((job-name (car job))
              (job-value (cdr job)))
            (if (or
                    (member job-name (cdr (assoc "after" job-value)))
                    (member job-name (cdr (assoc "before" job-value))))
                (raise-exception (make-job-depends-on-itself-exception job-name)))))
    (for-each job-depends-on-itself jobs)

    ; check all referenced jobs exist
    (define (all-referenced-jobs-exist jobs job-to-check)
        (define (check-list job-depender the-list)
            (for-each
                (lambda (job-name)
                    (if (not (assoc job-name jobs))
                        (raise-exception (make-job-not-found-exception job-depender job-name))))
                the-list))
        (let ((job-name (car job-to-check))
              (job-value (cdr job-to-check)))
            (check-list job-name (cdr (assoc "after" job-value)))
            (check-list job-name (cdr (assoc "before" job-value)))))
    (for-each (lambda (job-to-check) (all-referenced-jobs-exist jobs job-to-check)) jobs)

    ; check all targets exist
    (for-each
        (lambda (target-name)
            (if (not (assoc target-name jobs))
                (raise-exception (make-target-not-found-exception target-name))))
        targets)

    ; build list of jobs needed and check cyclic dependencies
    (define (enumerate-dependencies jobs job accumulator)
        (define (add-unique-jobs this-accumulator previous-accumulator)
            (lset-union string= previous-accumulator this-accumulator))
        (let* ((job-name job)
              (job-value (cdr (assoc job jobs)))
              (new-dependencies (cdr (assoc "after" job-value))))
            (cond
                ((null? new-dependencies) (list job-name))
                ((not (null? (lset-intersection string= accumulator new-dependencies)))
                    (raise-exception (make-cyclic-dependency-exception)))
                (else
                    (let ((updated-accumulator (append accumulator new-dependencies)))
                        (fold
                            (lambda (this previous)
                                (lset-union string= previous (enumerate-dependencies jobs this updated-accumulator)))
                            '()
                            new-dependencies))))))
    (define jobs-needed (append targets
        (fold
            (lambda (this previous)
                (lset-union string= previous (enumerate-dependencies jobs this '())))
            '()
            targets)))

    ; check all jobs are needed (issue warning if not)
    (define (check-job-is-needed jobs-needed job)
        (let ((job-name (car job)))
            (member job-name jobs-needed)))
    (define warnings
        (filter-map
            (lambda (job)
                (if (not (check-job-is-needed jobs-needed job)) (make-warning-job-not-needed (car job)) #f))
            jobs))

    (if (null? targets)
        (raise-exception (make-no-targets-exception)))

    (values jobs warnings))
