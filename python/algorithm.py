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

    # TODO: implementation

    return jobs_output, warnings
