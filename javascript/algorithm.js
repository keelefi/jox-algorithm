class AlgorithmError extends Error {
    constructor(message) {
        super(message);
        this.name = 'AlgorithmError';
    }
}

class CyclicDependencyError extends AlgorithmError {
    constructor() {
        super('Cyclic dependency detected');
        this.name = 'CyclicDependencyError';
    }
}

class JobDependsOnItselfError extends AlgorithmError {
    constructor(job_name) {
        super(`Job '${job_name}' depends on itself`);
        this.name = 'JobDependsOnItselfError';
        this.job_name = job_name;
    }
}

class NoTargetsError extends AlgorithmError {
    constructor(job_name) {
        super('No targets');
        this.name = 'NoTargetsError';
    }
}

class JobNotFoundError extends AlgorithmError {
    constructor(job_depender, job_dependee) {
        super(`Job '${job_depender}' references job '${job_dependee}', but job '${job_dependee}' does not exist`);
        this.name = 'JobNotFoundError';
        this.job_depender = job_depender;
        this.job_dependee = job_dependee;
    }
}

class TargetNotFoundError extends AlgorithmError {
    constructor(job_name) {
        super(`'targets' references job '${job_name}', but job '${job_name}' does not exist`);
        this.name = 'TargetNotFoundError';
        this.job_name = job_name;
    }
}

class AlgorithmWarning {
    constructor(message) {
        this.message = message;
    }
}

class JobNotRequiredWarning extends AlgorithmWarning {
    constructor(job_name) {
        super(`Job '${job_name}' is not required`);
        this.job_name = job_name;
    }
}

exports.algorithm = function (jobs, targets) {
    let jobs_actual = jobs;
    let warnings = [];

    // TODO: implementation

    return {
        jobs_actual: jobs_actual,
        warnings: warnings,
    };
}

exports.AlgorithmError = AlgorithmError;
exports.AlgorithmWarning = AlgorithmWarning;
