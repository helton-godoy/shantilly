#include <gtest/gtest.h>

#include "WidgetConfigs.h"

namespace {

using Sbx::Models::ActionConfig;
using Sbx::Models::ButtonConfig;
using Sbx::Models::WindowConfig;

TEST(WidgetConfigsTest, NamedControlsAreRequired) {
    ButtonConfig config;

    EXPECT_FALSE(config.isValid());

    config.name = "confirm";
    EXPECT_TRUE(config.isValid());
}

TEST(WidgetConfigsTest, WindowDimensionsMustBePositive) {
    WindowConfig config;

    EXPECT_TRUE(config.isValid());

    config.width = 0;
    EXPECT_FALSE(config.isValid());

    config.width = 800;
    config.height = -1;
    EXPECT_FALSE(config.isValid());
}

TEST(WidgetConfigsTest, ActionsValidateFieldsRequiredByTheirType) {
    ActionConfig action;

    action.type = ActionConfig::Shell;
    EXPECT_FALSE(action.isValid());
    action.command = "printf ready";
    EXPECT_TRUE(action.isValid());

    action = ActionConfig{};
    action.type = ActionConfig::Set;
    action.targetWidget = "status";
    EXPECT_FALSE(action.isValid());
    action.property = "text";
    EXPECT_TRUE(action.isValid());

    action = ActionConfig{};
    action.type = ActionConfig::Query;
    action.targetWidget = "username";
    EXPECT_FALSE(action.isValid());
    action.variable = "USER_NAME";
    EXPECT_TRUE(action.isValid());
}

} // namespace
