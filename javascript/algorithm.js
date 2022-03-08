const v8 = require('v8');
const structuredClone = obj => {
  return v8.deserialize(v8.serialize(obj));
};

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
    let jobs_output = structuredClone(jobs);
    let warnings = [];

    // check if targets is empty
    if (targets.length == 0) {
        throw new NoTargetsError();
    }

    // build all pairs
    for (const job in jobs_output) {
        if ('after' in jobs_output[job]) {
            for (const job_after of jobs_output[job]['after']) {
                if (!(job_after in jobs_output)) {
                    throw new JobNotFoundError(job, job_after);
                }
                if (job == job_after) {
                    throw new JobDependsOnItselfError(job);
                }
                if ('before' in jobs_output[job_after]) {
                    if (!(jobs_output[job_after]['before'].includes(job))) {
                        jobs_output[job_after]['before'].push(job);
                    }
                } else {
                    jobs_output[job_after]['before'] = [job];
                }
            }
        }

        if ('before' in jobs_output[job]) {
            for (const job_before of jobs_output[job]['before']) {
                if (!(job_before in jobs_output)) {
                    throw new JobNotFoundError(job, job_before);
                }
                if (job == job_before) {
                    throw new JobDependsOnItselfError(job);
                }
                if ('after' in jobs_output[job_before]) {
                    if (!(jobs_output[job_before]['after'].includes(job))) {
                        jobs_output[job_before]['after'].push(job);
                    }
                } else {
                    jobs_output[job_before]['after'] = [job];
                }
            }
        }
    }

    // collect list of needed jobs and check cyclic dependencies
    let jobs_needed = [];
    for (const target of targets) {
        if (!(target in jobs_output)) {
            throw new TargetNotFoundError(target);
        }

        let target_needs = [target];
        function check_cyclic_dependencies(jobs_output, jobs_needed, job_current, target_needs) {
            if ('after' in jobs_output[job_current]) {
                for (const job_after of jobs_output[job_current]['after']) {
                    if (jobs_needed.includes(job_after)) {
                        throw new CyclicDependencyError();
                    }
                    jobs_needed_copy = structuredClone(jobs_needed);
                    jobs_needed_copy.push(job_after);
                    if (!(job_after in target_needs)) {
                        target_needs.push(job_after);
                    }
                    check_cyclic_dependencies(jobs_output, jobs_needed_copy, job_after, target_needs);
                }
            }
        }
        check_cyclic_dependencies(jobs_output, [], target, target_needs);

        for (const job_needed of target_needs) {
            if (!(job_needed in jobs_needed)) {
                jobs_needed.push(job_needed);
            }
        }
    }

    // check that every job is needed
    for (const job in jobs) {
        if (!(jobs_needed.includes(job))) {
            warnings.push(new JobNotRequiredWarning(job));
        }
    }

    return {
        jobs_actual: jobs_output,
        warnings: warnings,
    };
}

exports.AlgorithmError = AlgorithmError;
exports.AlgorithmWarning = AlgorithmWarning;
