import copy

class AlgorithmException(Exception):
    def __init__(self, message):
        self.message = message
        super().__init__(self.message)

    def __eq__(self, other):
        return self.message == other.message

class CyclicDependencyException(AlgorithmException):
    def __init__(self):
        self.message = 'Cyclic dependency detected'
        super().__init__(self.message)

class JobDependsOnItselfException(AlgorithmException):
    def __init__(self, job_name):
        self.message = 'Job \'{}\' depends on itself'.format(job_name)
        super().__init__(self.message)

class NoTargetsException(AlgorithmException):
    def __init__(self):
        self.message = 'No targets'
        super().__init__(self.message)

class JobNotFoundException(AlgorithmException):
    def __init__(self, job_depender, job_dependee):
        message_template = 'Job \'{depender}\' references job \'{dependee}\', but job \'{dependee}\' does not exist'
        self.message = message_template.format(depender=job_depender, dependee=job_dependee)
        super().__init__(self.message)

class TargetNotFoundException(AlgorithmException):
    def __init__(self, job_name):
        message_template = '\'targets\' references job \'{job_name}\', but job \'{job_name}\' does not exist'
        self.message = message_template.format(job_name=job_name)
        super().__init__(self.message)

class AlgorithmWarning():
    def __init__(self, message):
        self.message = message

    def __eq__(self, other):
        return self.message == other.message

class JobNotRequiredWarning(AlgorithmWarning):
    def __init__(self, job_name):
        self.message = 'Job \'{}\' is not required'.format(job_name)
        super().__init__(self.message)

def algorithm(jobs, targets):
    jobs_output = copy.deepcopy(jobs)
    warnings = []

    if not targets:
        raise NoTargetsException()

    # build all pairs
    for job in jobs_output:
        if 'after' in jobs_output[job]:
            for job_after in jobs_output[job]['after']:
                if not job_after in jobs_output:
                    raise JobNotFoundException(job, job_after)
                if job == job_after:
                    raise JobDependsOnItselfException(job)
                if 'before' in jobs_output[job_after]:
                    if not job in jobs_output[job_after]['before']:
                        jobs_output[job_after]['before'].append(job)
                else:
                    jobs_output[job_after]['before'] = [job]
        if 'before' in jobs_output[job]:
            for job_before in jobs_output[job]['before']:
                if not job_before in jobs_output:
                    raise JobNotFoundException(job, job_before)
                if job == job_before:
                    raise JobDependsOnItselfException(job)
                if 'after' in jobs_output[job_before]:
                    if not job in jobs_output[job_before]['after']:
                        jobs_output[job_before]['after'].append(job)
                else:
                    jobs_output[job_before]['after'] = [job]

    # collect list of needed jobs and check cyclic dependencies
    jobs_needed = []
    for target in targets:
        if not target in jobs_output:
            raise TargetNotFoundException(target)

        target_needs = [target]

        def check_cyclic_dependencies(jobs_output, jobs_needed, job_current, target_needs):
            if 'after' in jobs_output[job_current]:
                for job_after in jobs_output[job_current]['after']:
                    if job_after in jobs_needed:
                        raise CyclicDependencyException()
                    jobs_needed_copy = copy.deepcopy(jobs_needed)
                    jobs_needed_copy.append(job_after)
                    if not job_after in target_needs:
                        target_needs.append(job_after)
                    check_cyclic_dependencies(jobs_output, jobs_needed_copy, job_after, target_needs)
        check_cyclic_dependencies(jobs_output, [], target, target_needs)

        for job_needed in target_needs:
            if not job_needed in jobs_needed:
                jobs_needed.append(job_needed)

    # check that every job is needed
    for job in jobs.keys():
        if not job in jobs_needed:
            warnings.append(JobNotRequiredWarning(job))

    return jobs_output, warnings
