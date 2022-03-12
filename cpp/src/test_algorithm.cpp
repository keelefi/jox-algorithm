#include "algorithm.h"
#include "exception.h"
#include "job.h"
#include "warning.h"

#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <nlohmann/json.hpp>

#include <map>
#include <filesystem>
#include <fstream>
#include <set>
#include <string>
#include <vector>

const std::filesystem::path PATH_TO_TESTS{"../tests/"};

const std::string ERROR_TYPE_WARNING{"WARNING"};
const std::string ERROR_TYPE_ERROR{"ERROR"};

std::ostream& operator<<(std::ostream& out, const algorithm::job& job_to_string)
{
    std::string result = "";

    auto after = job_to_string.get_after();
    auto before = job_to_string.get_before();

    if (!after.empty())
    {
        result += "after: [";
        for (const auto& str : after)
        {
            result += "'" + str + "', ";
        }
        result.erase(result.length()-2, 2);
        result += "]";
    }
    if (!before.empty())
    {
        if (!after.empty())
        {
            result += ", ";
        }
        result += "before: [";
        for (const auto& str : before)
        {
            result += "'" + str + "', ";
        }
        result.erase(result.length()-2, 2);
        result += "]";
    }

    out << result;
    return out;
}

std::ostream& operator<<(std::ostream& out, const algorithm::warning& warning_to_string)
{
    out << "Warning: {"
        << "type: " <<  warning_to_string.get_warning() << ", "
        << "message: \"" << warning_to_string.get_message() << "\"}";
    return out;
}

std::map<std::string, algorithm::job> jobs_from_json(const nlohmann::json& j)
{
    std::map<std::string, algorithm::job> result;

    for (auto& iter : j.items())
    {
        std::string name = iter.key();
        std::set<std::string> after;
        std::set<std::string> before;

        try
        {
            after = iter.value().at("after").get<std::set<std::string>>();
        }
        catch (nlohmann::json::out_of_range& e)
        {
        }

        try
        {
            before = iter.value().at("before").get<std::set<std::string>>();
        }
        catch (nlohmann::json::out_of_range& e)
        {
        }

        result.emplace(name, algorithm::job(after, before));
    }

    return result;
}

class AlgorithmTest: public testing::TestWithParam<std::filesystem::directory_entry>
{
};

std::vector<std::filesystem::directory_entry> test_vector_files()
{
    std::vector<std::filesystem::directory_entry> result;

    for (auto const& dir_entry : std::filesystem::directory_iterator{PATH_TO_TESTS})
    {
        if (!dir_entry.is_regular_file())
        {
            continue;
        }
        if (dir_entry.path().filename().string()[0] == '.')
        {
            continue;
        }
        if (dir_entry.path().extension() != ".json")
        {
            continue;
        }

        result.push_back(dir_entry);
    }

    return result;
}

TEST_P(AlgorithmTest, AlgorithmWorksCorrectly)
{
    // This might look a bit confusing.. The last parameter in INSTANTIATE_TEST_CASE_P is a ValueIn holding the result
    // from calling test_vector_files(). Each individual element from the vector is a new instance of this test case.
    // The individual element is retrieved with GetParam().
    auto dir_entry = GetParam();

    std::ifstream infilestream(dir_entry.path());
    nlohmann::json root = nlohmann::json::parse(infilestream);

    // read input_json into jobs
    auto jobs = jobs_from_json(root["input"]);

    // read targets
    auto targets = root.at("targets").get<std::vector<std::string>>();

    // read output_json into jobs_expected
    auto jobs_expected = jobs_from_json(root["output"]);

    // read warnings and errors
    std::set<algorithm::warning> warnings_expected;
    std::string exception_expected;
    for (const auto& error : root["errors"])
    {
        if (error["type"] == ERROR_TYPE_ERROR)
        {
            exception_expected = error["message"];
        }
        else if (error["type"] == ERROR_TYPE_WARNING)
        {
            warnings_expected.emplace(algorithm::warning_type_decode(error["error"]), error["message"]);
        }
        else
        {
            throw std::runtime_error("Invalid error type " + error["type"].get<std::string>());
        }
    }

    // Run the algorithm. If an exception is not expected, we just run the algorithm and check the results. If an
    // exception is expected, we try-catch it and check it was the correct one.
    if (exception_expected.empty())
    {
        auto warnings = algorithm::algorithm(jobs, targets);

        EXPECT_THAT(jobs, testing::ContainerEq(jobs_expected));
        EXPECT_EQ(warnings, warnings_expected);
    }
    else
    {
        // gtest has EXPECT_THROW, but it only checks the exception type. We want to also check the message.
        try
        {
            algorithm::algorithm(jobs, targets);

            FAIL() << "Expected exception: " << exception_expected;
        }
        catch (const algorithm::algorithm_exception& e)
        {
            EXPECT_EQ(exception_expected, e.what());
        }
    }
}

INSTANTIATE_TEST_CASE_P(
    AllTestCasesInstatiator,
    AlgorithmTest,
    testing::ValuesIn(test_vector_files()));

int main(int argc, char **argv)
{
    testing::InitGoogleTest(&argc, argv);
    return RUN_ALL_TESTS();
}

