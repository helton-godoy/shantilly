TEMPLATE = app
CONFIG += qt thread release c++17
QT += widgets svg charts

TARGET = shantilly
DESTDIR = bin
OBJECTS_DIR = obj
MOC_DIR = obj
RCC_DIR = obj
UI_DIR = obj

# Input
INCLUDEPATH += $$PWD/../../../libs/shantilly-ui/include
INCLUDEPATH += $$PWD/legacy/v2_incomplete
INCLUDEPATH += $$PWD/legacy/v1_monolith

HEADERS += SHantilly.h \
           legacy/v2_incomplete/tokenizer.h \
           legacy/v2_incomplete/command.h \
           legacy/v2_incomplete/command_registry.h \
           legacy/v2_incomplete/parser.h \
           legacy/v2_incomplete/parser_driver.h \
           legacy/v2_incomplete/execution_context.h \
           ../../../libs/shantilly-ui/include/icon_helper.h \
           ../../../libs/shantilly-ui/include/theme_manager.h \
           ../../../libs/shantilly-ui/include/logger.h \
           ../../../libs/shantilly-ui/include/custom_table_widget.h \
           ../../../libs/shantilly-ui/include/custom_chart_widget.h \
           legacy/v2_incomplete/commands/add_command.h \
           legacy/v2_incomplete/commands/set_command.h \
           legacy/v2_incomplete/commands/query_command.h \
           legacy/v2_incomplete/commands/simple_commands.h \
           legacy/v2_incomplete/commands/command_utils.h \
           legacy/v2_incomplete/commands/position_command.h \
           legacy/v2_incomplete/commands/unset_command.h

SOURCES += SHantilly.cc \
           legacy/v1_monolith/dialog_parser.cc \
           legacy/v1_monolith/dialog_main.cc \
           legacy/v1_monolith/dialog_set_options.cc \
           legacy/v1_monolith/dialog_slots.cc \
           legacy/v1_monolith/dialog_private.cc \
           legacy/v2_incomplete/tokenizer.cpp \
           legacy/v2_incomplete/command_registry.cpp \
           legacy/v2_incomplete/parser.cpp \
           legacy/v2_incomplete/parser_driver.cpp \
           legacy/v2_incomplete/execution_context.cpp \
           ../../../libs/shantilly-ui/src/icon_helper.cpp \
           ../../../libs/shantilly-ui/src/theme_manager.cpp \
           ../../../libs/shantilly-ui/src/logger.cpp \
           ../../../libs/shantilly-ui/src/custom_chart_widget.cpp \
           legacy/v2_incomplete/commands/add_command.cpp \
           legacy/v2_incomplete/commands/set_command.cpp \
           legacy/v2_incomplete/commands/query_command.cpp \
           legacy/v2_incomplete/commands/simple_commands.cpp \
           legacy/v2_incomplete/commands/command_utils.cpp \
           legacy/v2_incomplete/commands/position_command.cpp \
           legacy/v2_incomplete/commands/unset_command.cpp

# install recipe options
target.path = /usr/bin
INSTALLS += target

# Added C/C++ compiler options
# Strict warnings
QMAKE_CXXFLAGS += -Werror -Wall -Wextra

# Windows specific: Ensure console is available for stdout/stdin
win32:CONFIG += console