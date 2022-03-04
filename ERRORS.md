# Errors

This section lists the possible errors the algorithm can encounter. Note that the algorithm needs to report each error
it encounters, respectively.

## Cyclic Dependency

**Description**: Jobs depend on each other in a cyclic manner and therefore no execution order can be found.

**Error**: `CYCLIC DEPENDENCY`

**Type**: `ERROR`

**Example**: Job `B` is `after` job `A` and job `A` is `after` job `B`.

## Job Depends On Itself

**Description**: Job is defined to be `after` or `before` itself.

**Error**: `JOB DEPENDS ON ITSELF`

**Type**: `ERROR`

**Example**: Job `A` is `after` job `A`.

## No Targets

**Description**: No targets have been defined.

**Error**: `NO TARGETS`

**Type**: `ERROR`

**Example**: The list of targets is empty.

## Job Not Found

**Description**: A job references a job that does not exist.

**Error**: `JOB NOT FOUND`

**Type**: `ERROR`

**Example**: Job `A` is `after` job `foobar`, but job `foobar` has not been defined anywhere.

## Job Not Required

**Description**: Job is not required.

**Error**: `JOB NOT REQUIRED`

**Type**: `WARNING`

**Example**: Job `B` is the only `target` and does not depend on job `A`.
