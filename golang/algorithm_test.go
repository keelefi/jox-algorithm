package algorithm

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "testing"

    "github.com/google/go-cmp/cmp"
)

const TEST_FILES_PATH = "../tests"

type TestCaseData struct {
    input           map[string]Job
    targets         map[string]bool
    output          map[string]Job
    warnings        []AlgorithmWarning
    errorExpected   AlgorithmError
}

func readStringSet(data interface{}, key string) (result map[string]bool) {
    valueInterface, ok := data.(map[string]interface{})[key]
    if !ok {
        return make(map[string]bool, 0)
    }

    valueInterfaceSlice := valueInterface.([]interface{})
    result = make(map[string]bool, len(valueInterfaceSlice))
    for _, v := range valueInterfaceSlice {
        result[fmt.Sprintf("%s", v)] = true
    }

    return result
}

func jsonToJob(key string, value interface{}) (result Job) {
    result.After = readStringSet(value, "after")
    result.Before = readStringSet(value, "before")

    return result
}

func filterFiles(files []os.FileInfo) (result []os.FileInfo) {
    for _, v := range files {
        if v.IsDir() {
            continue
        }
        if strings.HasPrefix(v.Name(), ".") {
            continue
        }
        if !strings.HasSuffix(v.Name(), ".json") {
            continue
        }
        result = append(result, v)
    }

    return result
}

func getTestFiles(t *testing.T) (result []os.FileInfo, err error) {
    files, err := ioutil.ReadDir(TEST_FILES_PATH)
    if err != nil {
        t.Fatalf(fmt.Sprintf("Could not read test files path: %s", TEST_FILES_PATH))
        return nil, err
    }

    result = filterFiles(files)

    return result, nil
}

func parseJobs(t *testing.T, data map[string]interface{}) (result map[string]Job) {
    result = make(map[string]Job, len(data))

    for key, value := range data {
        result[key] = jsonToJob(key, value)
    }

    return result
}

func parseTargets(t *testing.T, data []interface{}) (result map[string]bool) {
    result = make(map[string]bool, len(data))

    for _, v := range data {
        result[fmt.Sprintf("%s", v)] = true
    }

    return result
}

func parseErrors(t *testing.T, data []interface{}) (warnings []AlgorithmWarning, errorExpected AlgorithmError) {
    for _, v := range data {
        v_map := v.(map[string]interface{})
        errorType, err := GetErrorType(v_map["type"].(string))
        if err != nil {
            t.Fatalf("error: %s", err)
        }
        if errorType == ErrorTypeError {
            if errorExpected.errorEnumeration != ErrorNone {
                t.Fatalf("Error set twice")
            }
            errorEnumeration, err := GetErrorEnumeration(v_map["error"].(string))
            if err != nil {
                t.Fatalf("error: %s", err)
            }
            errorExpected.errorEnumeration = errorEnumeration
            errorExpected.message = v_map["message"].(string)
        } else if errorType == ErrorTypeWarning {
            warningEnumeration, err := GetWarningEnumeration(v_map["error"].(string))
            if err != nil {
                t.Fatalf("error: %s", err)
            }
            newWarning := AlgorithmWarning{
                Enumeration:    warningEnumeration,
                Message:        v_map["message"].(string),
            }
            warnings = append(warnings, newWarning)
        }
    }

    return warnings, errorExpected
}

func parseFile(t *testing.T, filename string) (result TestCaseData) {
    content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", TEST_FILES_PATH, filename))
    if err != nil {
        t.Fatalf("Could not read JSON file: %s", filename)
    }

    var data map[string]interface{}
    err = json.Unmarshal(content, &data)
    if err != nil {
        t.Fatalf("Unmarshalling %s failed", filename)
    }

    result.input = parseJobs(t, data["input"].(map[string]interface{}))
    result.output = parseJobs(t, data["output"].(map[string]interface{}))
    result.targets = parseTargets(t, data["targets"].([]interface{}))
    result.warnings, result.errorExpected = parseErrors(t, data["errors"].([]interface{}))

    return result
}

func TestAlgorithm(t *testing.T) {
    testFiles, err := getTestFiles(t)
    if err != nil {
        return
    }

    for _, testFile := range testFiles {
        t.Run(testFile.Name(), func (t *testing.T) {
            testCaseData := parseFile(t, testFile.Name())

            jobsActual, warningsActual, err := Algorithm(testCaseData.input, testCaseData.targets)

            if testCaseData.errorExpected.errorEnumeration == ErrorNone {
                if err != nil {
                    t.Errorf("did not expect error, got: '%s'", err.Error())
                }

                // No error expected, check jobsActual and warningsActual
                jobsDiff := cmp.Diff(jobsActual, testCaseData.output)
                if jobsDiff != "" {
                    t.Errorf("jobsActual does not match jobsExpected:\n%s", jobsDiff)
                }
                warningsDiff := cmp.Diff(warningsActual, testCaseData.warnings)
                if warningsDiff != "" {
                    t.Errorf("warningsActual does not match warningsExpected:\n%s", warningsDiff)
                }
            } else {
                // Error expected, ignore other values and only check err
                if err != nil {
                    if err.Error() != testCaseData.errorExpected.message {
                        t.Errorf("Expected error '%s', but got '%s'", testCaseData.errorExpected.message, err.Error())
                    }
                } else {
                    t.Errorf("Expected error '%s', but got none", testCaseData.errorExpected.message)
                }
            }
        })
    }

    return
}

