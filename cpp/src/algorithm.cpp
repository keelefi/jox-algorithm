#include "algorithm.h"

#include "exception.h"

#include <functional>

namespace algorithm
{

std::set<warning> algorithm(std::map<std::string, job>& jobs, const std::vector<std::string>& targets)
{
    std::set<warning> warnings;

    if (targets.empty())
    {
        throw no_targets_exception();
    }

    // build all pairs
    for (auto& job_iter : jobs)
    {
        for (auto& job_after_str : job_iter.second.get_after())
        {
            auto job_after = jobs.find(job_after_str);
            if (job_after == jobs.end())
            {
                throw job_not_found_exception(job_iter.first, job_after_str);
            }
            if (job_iter.first == job_after_str)
            {
                throw job_depends_on_itself_exception(job_iter.first);
            }
            auto jobs_before = job_after->second.get_before();
            auto job_after_before = jobs_before.find(job_iter.first);
            if (job_after_before == jobs_before.end())
            {
                job_after->second.add_before(job_iter.first);
            }
        }
        for (auto& job_before_str : job_iter.second.get_before())
        {
            auto job_before = jobs.find(job_before_str);
            if (job_before == jobs.end())
            {
                throw job_not_found_exception(job_iter.first, job_before_str);
            }
            if (job_iter.first == job_before_str)
            {
                throw job_depends_on_itself_exception(job_iter.first);
            }
            auto jobs_after = job_before->second.get_before();
            auto job_before_after = jobs_after.find(job_iter.first);
            if (job_before_after == jobs_after.end())
            {
                job_before->second.add_after(job_iter.first);
            }
        }
    }

    // collect list of needed jobs and check cyclic dependencies
    std::set<std::string> jobs_needed;
    for (const auto& target : targets)
    {
        auto target_job = jobs.find(target);
        if (target_job == jobs.end())
        {
            throw target_not_found_exception(target);
        }

        std::set<std::string> target_needs = {target};

        std::function<void(const std::set<std::string>, const std::string)>
        check_cyclic_dependencies = [&jobs, &target_needs, &check_cyclic_dependencies](
                const std::set<std::string>& jobs_needed,
                const std::string& job_current_str) -> void
        {
            const auto job_current = jobs.find(job_current_str);
            for (const auto& job_after : job_current->second.get_after())
            {
                if (jobs_needed.find(job_after) != jobs_needed.end())
                {
                    throw cyclic_dependency_exception();
                }
                target_needs.insert(job_after);
                auto jobs_needed_copy(jobs_needed);
                jobs_needed_copy.insert(job_after);
                check_cyclic_dependencies(jobs_needed_copy, job_after);
            }
        };
        check_cyclic_dependencies({}, target);

        for (const auto& job_needed : target_needs)
        {
            jobs_needed.insert(job_needed);
        }
    }

    // check that every job is needed
    for (const auto& job_iter : jobs)
    {
        if (jobs_needed.find(job_iter.first) == jobs_needed.end())
        {
            warnings.emplace(job_not_required_warning(job_iter.first));
        }
    }

    return warnings;
}

} // namespace algorithm

