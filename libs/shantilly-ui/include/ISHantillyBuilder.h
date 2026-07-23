#ifndef ISHANTILLYBUILDER_H
#define ISHANTILLYBUILDER_H

#include <QLayout>
#include <QWidget>

#include "WidgetConfigs.h"

// Forward declaration if needed for legacy support
class PushButtonWidget;

class ISHantillyBuilder {
public:
    virtual ~ISHantillyBuilder() = default;

    // --- New Architecture Methods (Phase 1 & 2) ---

    virtual QWidget* buildWindow(const Sbx::Models::WindowConfig& config) = 0;
    virtual QWidget* buildButton(const Sbx::Models::ButtonConfig& config) = 0;
    virtual QWidget* buildLabel(const Sbx::Models::LabelConfig& config) = 0;
    virtual QWidget* buildLineEdit(const Sbx::Models::LineEditConfig& config) = 0;
    virtual QWidget* buildComboBox(const Sbx::Models::ComboBoxConfig& config) = 0;
    virtual QWidget* buildList(const Sbx::Models::ListConfig& config) = 0;
    virtual QWidget* buildTable(const Sbx::Models::TableConfig& config) = 0;
    virtual QWidget* buildProgressBar(const Sbx::Models::ProgressBarConfig& config) = 0;
    virtual QWidget* buildChart(const Sbx::Models::ChartConfig& config) = 0;
    virtual QWidget* buildCheckBox(const Sbx::Models::CheckBoxConfig& config) = 0;
    virtual QWidget* buildRadioButton(const Sbx::Models::RadioButtonConfig& config) = 0;
    virtual QWidget* buildCalendar(const Sbx::Models::CalendarConfig& config) = 0;
    virtual QWidget* buildSeparator(const Sbx::Models::SeparatorConfig& config) = 0;
    virtual QWidget* buildSpinBox(const Sbx::Models::SpinBoxConfig& config) = 0;
    virtual QWidget* buildSlider(const Sbx::Models::SliderConfig& config) = 0;
    virtual QWidget* buildTextEdit(const Sbx::Models::TextEditConfig& config) = 0;
    virtual QWidget* buildGroupBox(const Sbx::Models::GroupBoxConfig& config) = 0;
    virtual QWidget* buildFrame(const Sbx::Models::FrameConfig& config) = 0;
    virtual QWidget* buildTabWidget(const Sbx::Models::TabWidgetConfig& config) = 0;
    virtual QLayout* buildLayoutStructure(const Sbx::Models::SbxLayoutConfig& config) = 0;

    // --- Legacy Methods (To be deprecated/removed) ---
    virtual PushButtonWidget* buildPushButton(const QString& title, const QString& name) = 0;
};

#endif // ISHANTILLYBUILDER_H