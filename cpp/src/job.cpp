#include "job.h"

namespace algorithm
{

job::job(const std::set<std::string>& after, const std::set<std::string>& before) :
        m_after(after),
        m_before(before)
{
}

bool job::operator==(const job& other) const
{
    return (m_after == other.m_after) &&
            (m_before == other.m_before);
}

void job::add_after(const std::string& job_after)
{
    m_after.insert(job_after);
}

void job::add_before(const std::string& job_before)
{
    m_before.insert(job_before);
}

const std::set<std::string>& job::get_after() const
{
    return m_after;
}

const std::set<std::string>& job::get_before() const
{
    return m_before;
}

} // namespace algorithm

