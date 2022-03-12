#pragma once

#include <map>
#include <set>
#include <string>
#include <vector>

namespace algorithm
{

class job
{
public:
    job(const std::set<std::string>&, const std::set<std::string>&);

    // TODO: default comparison is available in C++20
    //bool operator==(const job&) const = default;
    bool operator==(const job&) const;

    void add_after(const std::string&);
    void add_before(const std::string&);

    const std::set<std::string>& get_after() const;
    const std::set<std::string>& get_before() const;

private:
    std::set<std::string> m_after;
    std::set<std::string> m_before;
};

} // namespace algorithm

