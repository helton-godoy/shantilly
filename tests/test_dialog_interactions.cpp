/*
 * GUI widgets for shell scripts - SHantilly version 1.0
 *
 * Copyright (C) 2015-2016, 2020 Andriy Martynets <andy.martynets@gmail.com>
 *------------------------------------------------------------------------------
 * This file is part of SHantilly.
 *
 * SHantilly is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * SHantilly is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with SHantilly. If not, see http://www.gnu.org/licenses/.
 *------------------------------------------------------------------------------
 */

#include <gtest/gtest.h>

#include <QPushButton>
#include <QTest>
#include <QTimer>
#include <cstdio>

#include "SHantilly.h"

namespace {

class OutputCapture {
public:
    OutputCapture() : file_(std::tmpfile()) {}
    ~OutputCapture() {
        if (file_)
            std::fclose(file_);
    }

    FILE* file() const { return file_; }

    QByteArray readAll() {
        if (!file_)
            return {};

        std::fflush(file_);
        std::rewind(file_);

        QByteArray output;
        char buffer[256];
        std::size_t bytesRead;
        while ((bytesRead = std::fread(buffer, 1, sizeof(buffer), file_)) > 0)
            output.append(buffer, static_cast<qsizetype>(bytesRead));
        return output;
    }

private:
    FILE* file_;
};

QPushButton* requireButton(SHantilly& dialog, const char* name) {
    auto* button = dialog.findChild<QPushButton*>(QString::fromUtf8(name));
    EXPECT_NE(button, nullptr);
    return button;
}

TEST(DialogInteractionsTest, ApplyReportsValuesWithoutClosing) {
    OutputCapture output;
    ASSERT_NE(output.file(), nullptr);

    SHantilly dialog("Test", nullptr, false, output.file());
    dialog.addCheckBox("Option", "option", true);
    dialog.addPushButton("Apply", "apply", true, false);
    dialog.show();

    auto* button = requireButton(dialog, "apply");
    ASSERT_NE(button, nullptr);
    QTest::mouseClick(button, Qt::LeftButton);

    EXPECT_TRUE(dialog.isVisible());
    EXPECT_EQ(output.readAll(), QByteArray("apply=clicked\noption=1\n"));
}

TEST(DialogInteractionsTest, AcceptReturnsOneAndReportsValues) {
    OutputCapture output;
    ASSERT_NE(output.file(), nullptr);

    SHantilly dialog("Test", nullptr, false, output.file());
    dialog.addCheckBox("Option", "option", true);
    dialog.addPushButton("Confirm", "confirm", true, true);
    dialog.show();

    auto* button = requireButton(dialog, "confirm");
    ASSERT_NE(button, nullptr);
    QTimer::singleShot(0, button, &QPushButton::click);

    EXPECT_EQ(QApplication::exec(), 1);
    EXPECT_EQ(output.readAll(), QByteArray("confirm=clicked\noption=1\n"));
}

TEST(DialogInteractionsTest, RejectReturnsZeroWithoutReportingValues) {
    OutputCapture output;
    ASSERT_NE(output.file(), nullptr);

    SHantilly dialog("Test", nullptr, false, output.file());
    dialog.addCheckBox("Option", "option", true);
    dialog.addPushButton("Cancel", "cancel", false, true);
    dialog.show();

    auto* button = requireButton(dialog, "cancel");
    ASSERT_NE(button, nullptr);
    QTimer::singleShot(0, button, &QPushButton::click);

    EXPECT_EQ(QApplication::exec(), 0);
    EXPECT_EQ(output.readAll(), QByteArray("cancel=clicked\n"));
}

} // namespace
