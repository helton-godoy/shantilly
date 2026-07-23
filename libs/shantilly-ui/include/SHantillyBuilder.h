#ifndef SHANTILLYBUILDER_H
#define SHANTILLYBUILDER_H

#include <QObject>

#include "ISHantillyBuilder.h"

class SHantillyBuilder : public QObject, public ISHantillyBuilder {
    Q_OBJECT
public:
    explicit SHantillyBuilder(QObject* parent = nullptr);

    QWidget* buildWindow(const Sbx::Models::WindowConfig& config) override;
    QWidget* buildButton(const Sbx::Models::ButtonConfig& config) override;
    QWidget* buildLabel(const Sbx::Models::LabelConfig& config) override;
    QWidget* buildLineEdit(const Sbx::Models::LineEditConfig& config) override;
    QWidget* buildComboBox(const Sbx::Models::ComboBoxConfig& config) override;
    QWidget* buildList(const Sbx::Models::ListConfig& config) override;
    QWidget* buildTable(const Sbx::Models::TableConfig& config) override;
    QWidget* buildProgressBar(const Sbx::Models::ProgressBarConfig& config) override;
    QWidget* buildChart(const Sbx::Models::ChartConfig& config) override;
    QWidget* buildCheckBox(const Sbx::Models::CheckBoxConfig& config) override;
    QWidget* buildRadioButton(const Sbx::Models::RadioButtonConfig& config) override;
    QWidget* buildCalendar(const Sbx::Models::CalendarConfig& config) override;
    QWidget* buildSeparator(const Sbx::Models::SeparatorConfig& config) override;
    QWidget* buildSpinBox(const Sbx::Models::SpinBoxConfig& config) override;
    QWidget* buildSlider(const Sbx::Models::SliderConfig& config) override;
    QWidget* buildTextEdit(const Sbx::Models::TextEditConfig& config) override;
    QWidget* buildGroupBox(const Sbx::Models::GroupBoxConfig& config) override;
    QWidget* buildFrame(const Sbx::Models::FrameConfig& config) override;
    QWidget* buildTabWidget(const Sbx::Models::TabWidgetConfig& config) override;
    QLayout* buildLayoutStructure(const Sbx::Models::SbxLayoutConfig& config) override;

    // Legacy support
    PushButtonWidget* buildPushButton(const QString& title, const QString& name) override;
};

#endif // SHANTILLYBUILDER_H
