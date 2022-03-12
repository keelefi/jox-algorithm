#pragma once

#include <map>
#include <set>
#include <string>
#include <vector>

#include "job.h"
#include "warning.h"

namespace algorithm
{

std::set<warning> algorithm(std::map<std::string, job>& jobs, const std::vector<std::string>& targets);

} // namespace algorithm

