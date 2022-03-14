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
    warningEnumeration  WarningEnumeration
    message             string
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

func Algorithm(input map[string]Job, targets map[string]bool) (map[string]Job, []AlgorithmWarning, error) {
    output := DeepCopyJobs(input)
    var warnings []AlgorithmWarning

    // TODO: implementation

    return output, warnings, nil
}

