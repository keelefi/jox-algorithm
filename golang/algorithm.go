package algorithm

import (
    "errors"
    "fmt"
)

type Job struct {
    After   map[string]bool
    Before  map[string]bool
}

func SetToString(set map[string]bool) (result string) {
    result = "["

    for k, _ := range set {
        result += fmt.Sprintf(" %s ", k)
    }

    result += "]"

    return result
}

func (j Job) Copy() (result Job) {
    result.After = make(map[string]bool, len(j.After))
    for k, v := range j.After {
        result.After[k] = v
    }

    result.Before = make(map[string]bool, len(j.Before))
    for k, v := range j.Before {
        result.Before[k] = v
    }

    return result
}

func DeepCopyJobs(jobs map[string]Job) (result map[string]Job) {
    result = make(map[string]Job, len(jobs))

    for i, v := range jobs {
        result[i] = v.Copy()
    }

    return result
}

type ErrorType int

const (
    ErrorTypeNone ErrorType = iota
    ErrorTypeError
    ErrorTypeWarning
)

type ErrorEnumeration int

const (
    ErrorNone ErrorEnumeration = iota
    ErrorCyclicDependency
    ErrorJobDependsOnItself
    ErrorNoTargets
    ErrorJobNotFound
)

type WarningEnumeration int

const (
    WarningJobNotRequired = iota
)

func GetErrorType(errorType string) (ErrorType, error) {
    if errorType == "ERROR" {
        return ErrorTypeError, nil
    }
    if errorType == "WARNING" {
        return ErrorTypeWarning, nil
    }
    return ErrorTypeNone, errors.New(fmt.Sprintf("Unknown error type: %s", errorType))
}

func (e ErrorEnumeration) String() string {
    switch e {
    case ErrorCyclicDependency:
        return "CYCLIC DEPENDENCY"
    case ErrorJobDependsOnItself:
        return "JOB DEPENDS ON ITSELF"
    case ErrorNoTargets:
        return "NO TARGETS"
    case ErrorJobNotFound:
        return "JOB NOT FOUND"
    }
    return "ERROR"
}

func GetErrorEnumeration(errorString string) (ErrorEnumeration, error) {
    if errorString == "CYCLIC DEPENDENCY" {
        return ErrorCyclicDependency, nil
    }
    if errorString == "JOB DEPENDS ON ITSELF" {
        return ErrorJobDependsOnItself, nil
    }
    if errorString == "NO TARGETS" {
        return ErrorNoTargets, nil
    }
    if errorString == "JOB NOT FOUND" {
        return ErrorJobNotFound, nil
    }
    return 0, errors.New(fmt.Sprintf("Unknown error enumeration: %s", errorString))
}

func (w WarningEnumeration) String() string {
    switch w {
    case WarningJobNotRequired:
        return "JOB NOT REQUIRED"
    }
    return "ERROR"
}

func GetWarningEnumeration(warningString string) (WarningEnumeration, error) {
    if warningString == "JOB NOT REQUIRED" {
        return WarningJobNotRequired, nil
    }
    return 0, errors.New(fmt.Sprintf("Unknown warning enumeration: %s", warningString))
}

type AlgorithmError struct {
    errorEnumeration    ErrorEnumeration
    message             string
}

type AlgorithmWarning struct {
    Enumeration WarningEnumeration
    Message     string
}

type CyclicDependencyError struct {}
type JobDependsOnItselfError struct {
    jobName string
}
type NoTargetsError struct {}
type JobNotFoundError struct {
    depender    string
    dependee    string
}
type TargetNotFoundError struct {
    jobName string
}

type JobNotRequiredWarning struct {
    jobName string
}

func (e *CyclicDependencyError) Error() string {
    return "Cyclic dependency detected"
}

func (e *JobDependsOnItselfError) Error() string {
    return fmt.Sprintf("Job '%s' depends on itself", e.jobName)
}

func (e *NoTargetsError) Error() string {
    return "No targets"
}

func (e *JobNotFoundError) Error() string {
    return fmt.Sprintf("Job '%s' references job '%s', but job '%s' does not exist", e.depender, e.dependee, e.dependee)
}

func (e *TargetNotFoundError) Error() string {
    return fmt.Sprintf("'targets' references job '%s', but job '%s' does not exist", e.jobName, e.jobName)
}

func (w *JobNotRequiredWarning) Warning() string {
    return fmt.Sprintf("Job '%s' is not required", w.jobName)
}

func (w *JobNotRequiredWarning) ToWarning() AlgorithmWarning {
    return AlgorithmWarning {
        Enumeration:    WarningJobNotRequired,
        Message:        w.Warning(),
    }
}

func Algorithm(input map[string]Job, targets map[string]bool) (map[string]Job, []AlgorithmWarning, error) {
    output := DeepCopyJobs(input)
    var warnings []AlgorithmWarning

    if len(targets) == 0 {
        return nil, nil, &NoTargetsError{}
    }

    // build all pairs
    for k, v := range output {
        for jobAfterStr := range v.After {
            jobAfter, ok := output[jobAfterStr]
            if !ok {
                err := &JobNotFoundError{
                    depender: k,
                    dependee: jobAfterStr,
                }
                return nil, nil, err
            }
            if jobAfterStr == k {
                err := &JobDependsOnItselfError{
                    jobName: k,
                }
                return nil, nil, err
            }
            if _, ok := jobAfter.Before[k]; !ok {
                jobAfter.Before[k] = true
            }
        }
        for jobBeforeStr := range v.Before {
            jobBefore, ok := output[jobBeforeStr]
            if !ok {
                err := &JobNotFoundError{
                    depender: k,
                    dependee: jobBeforeStr,
                }
                return nil, nil, err
            }
            if jobBeforeStr == k {
                err := &JobDependsOnItselfError{
                    jobName: k,
                }
                return nil, nil, err
            }
            if _, ok := jobBefore.After[k]; !ok {
                jobBefore.After[k] = true
            }
        }
    }

    // collect list of needed jobs and check cyclic dependencies
    jobsNeeded := make(map[string]bool)
    for target, _ := range targets {
        if _, ok := output[target]; !ok {
            err := &TargetNotFoundError{
                jobName: target,
            }
            return nil, nil, err
        }

        targetNeeds := make(map[string]bool)
        targetNeeds[target] = true

        var checkCyclicDependencies func(map[string]bool, string) error
        checkCyclicDependencies = func(jobsNeeded map[string]bool, jobCurrentStr string) error {
            jobCurrent := output[jobCurrentStr]
            for jobAfterStr, _ := range jobCurrent.After {
                if _, ok := jobsNeeded[jobAfterStr]; ok {
                    return &CyclicDependencyError{}
                }
                targetNeeds[jobAfterStr] = true
                jobsNeededCopy := make(map[string]bool, len(jobsNeeded)+1)
                for k, _ := range jobsNeeded {
                    jobsNeededCopy[k] = true
                }
                jobsNeededCopy[jobCurrentStr] = true
                if err := checkCyclicDependencies(jobsNeededCopy, jobAfterStr); err != nil {
                    return err
                }
            }
            return nil
        }
        if err := checkCyclicDependencies(make(map[string]bool), target); err != nil {
            return nil, nil, err
        }

        for k, _ := range targetNeeds {
            jobsNeeded[k] = true
        }
    }

    // check that every job is needed
    for k, _ := range output {
        if _, ok := jobsNeeded[k]; !ok {
            w := JobNotRequiredWarning{jobName: k}
            warnings = append(warnings, w.ToWarning())
        }
    }

    return output, warnings, nil
}

