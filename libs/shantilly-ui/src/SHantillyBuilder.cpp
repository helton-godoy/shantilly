#include "SHantillyBuilder.h"

#include <QCalendarWidget>
#include <QCheckBox>
#include <QComboBox>
#include <QDialog>
#include <QFrame>
#include <QGroupBox>
#include <QLabel>
#include <QLineEdit>
#include <QListWidget>
#include <QProgressBar>
#include <QPushButton>
#include <QRadioButton>
#include <QSlider>
#include <QSpinBox>
#include <QTabWidget>
#include <QTableWidget>
#include <QTextEdit>
#include <QVBoxLayout>

#include "custom_chart_widget.h"
#include "push_button_widget.h"

SHantillyBuilder::SHantillyBuilder(QObject* parent) : QObject(parent) {
}

QWidget* SHantillyBuilder::buildWindow(const Sbx::Models::WindowConfig& config) {
    auto* window = new QDialog();
    window->setWindowTitle(config.title);
    window->resize(config.width, config.height);

    // Default layout for new windows to allow immediate adding of widgets
    if (!window->layout()) {
        new QVBoxLayout(window);
    }

    return window;
}

QWidget* SHantillyBuilder::buildButton(const Sbx::Models::ButtonConfig& config) {
    // For now using standard QPushButton, later can use PushButtonWidget if it adds value
    auto* btn = new QPushButton(config.text);
    btn->setObjectName(config.name);
    btn->setCheckable(config.checkable);
    btn->setChecked(config.checked);
    return btn;
}

QWidget* SHantillyBuilder::buildLabel(const Sbx::Models::LabelConfig& config) {
    auto* lbl = new QLabel(config.text);
    lbl->setObjectName(config.name);
    lbl->setWordWrap(config.wordWrap);
    return lbl;
}

QWidget* SHantillyBuilder::buildLineEdit(const Sbx::Models::LineEditConfig& config) {
    auto* le = new QLineEdit(config.text);
    le->setObjectName(config.name);
    le->setPlaceholderText(config.placeholder);
    if (config.passwordMode) {
        le->setEchoMode(QLineEdit::Password);
    }
    return le;
}

QWidget* SHantillyBuilder::buildComboBox(const Sbx::Models::ComboBoxConfig& config) {
    auto* cb = new QComboBox();
    cb->setObjectName(config.name);
    cb->addItems(config.items);
    if (config.currentIndex >= 0 && config.currentIndex < cb->count()) {
        cb->setCurrentIndex(config.currentIndex);
    }
    return cb;
}

QWidget* SHantillyBuilder::buildList(const Sbx::Models::ListConfig& config) {
    auto* lw = new QListWidget();
    lw->setObjectName(config.name);
    lw->addItems(config.items);
    if (config.multipleSelection) {
        lw->setSelectionMode(QAbstractItemView::MultiSelection);
    }
    return lw;
}

QWidget* SHantillyBuilder::buildTable(const Sbx::Models::TableConfig& config) {
    auto* tw = new QTableWidget();
    tw->setObjectName(config.name);
    tw->setColumnCount(config.headers.size());
    tw->setHorizontalHeaderLabels(config.headers);
    tw->setRowCount(config.rows.size());
    for (int r = 0; r < config.rows.size(); ++r) {
        for (int c = 0; c < config.rows[r].size() && c < config.headers.size(); ++c) {
            tw->setItem(r, c, new QTableWidgetItem(config.rows[r][c]));
        }
    }
    return tw;
}

QWidget* SHantillyBuilder::buildProgressBar(const Sbx::Models::ProgressBarConfig& config) {
    auto* pb = new QProgressBar();
    pb->setObjectName(config.name);
    pb->setMinimum(config.minimum);
    pb->setMaximum(config.maximum);
    pb->setValue(config.value);
    pb->setFormat(config.format);
    return pb;
}

QWidget* SHantillyBuilder::buildChart(const Sbx::Models::ChartConfig& config) {
    auto* chart = new CustomChartWidget();
    chart->setObjectName(config.name);
    chart->setChartTitle(config.title);

    // Populate chart data from config
    for (auto it = config.data.begin(); it != config.data.end(); ++it) {
        chart->addPoint(it.key(), it.value());
    }

    return chart;
}

QWidget* SHantillyBuilder::buildCheckBox(const Sbx::Models::CheckBoxConfig& config) {
    auto* cb = new QCheckBox(config.text);
    cb->setObjectName(config.name);
    cb->setChecked(config.checked);
    return cb;
}

QWidget* SHantillyBuilder::buildRadioButton(const Sbx::Models::RadioButtonConfig& config) {
    auto* rb = new QRadioButton(config.text);
    rb->setObjectName(config.name);
    rb->setChecked(config.checked);
    return rb;
}

QWidget* SHantillyBuilder::buildCalendar(const Sbx::Models::CalendarConfig& config) {
    auto* cw = new QCalendarWidget();
    cw->setObjectName(config.name);
    return cw;
}

QWidget* SHantillyBuilder::buildSeparator(const Sbx::Models::SeparatorConfig& config) {
    auto* line = new QFrame();
    line->setObjectName(config.name);
    if (config.orientation == Qt::Horizontal) {
        line->setFrameShape(QFrame::HLine);
    } else {
        line->setFrameShape(QFrame::VLine);
    }
    line->setFrameShadow(QFrame::Sunken);
    return line;
}

QWidget* SHantillyBuilder::buildSpinBox(const Sbx::Models::SpinBoxConfig& config) {
    auto* sb = new QSpinBox();
    sb->setObjectName(config.name);
    sb->setMinimum(config.min);
    sb->setMaximum(config.max);
    sb->setValue(config.value);
    sb->setSingleStep(config.step);
    sb->setPrefix(config.prefix);
    sb->setSuffix(config.suffix);
    return sb;
}

QWidget* SHantillyBuilder::buildSlider(const Sbx::Models::SliderConfig& config) {
    auto* sl = new QSlider();
    sl->setObjectName(config.name);
    sl->setMinimum(config.min);
    sl->setMaximum(config.max);
    sl->setValue(config.value);
    sl->setOrientation(static_cast<Qt::Orientation>(config.orientation));
    return sl;
}

QWidget* SHantillyBuilder::buildTextEdit(const Sbx::Models::TextEditConfig& config) {
    auto* te = new QTextEdit();
    te->setObjectName(config.name);
    if (config.richText) {
        te->setHtml(config.text);
    } else {
        te->setPlainText(config.text);
    }
    te->setReadOnly(config.readOnly);
    return te;
}

QWidget* SHantillyBuilder::buildGroupBox(const Sbx::Models::GroupBoxConfig& config) {
    auto* gb = new QGroupBox(config.title);
    gb->setObjectName(config.name);

    // Build and set layout if valid
    QLayout* layout = buildLayoutStructure(config.layout);
    if (layout) {
        gb->setLayout(layout);
    }

    return gb;
}

QWidget* SHantillyBuilder::buildFrame(const Sbx::Models::FrameConfig& config) {
    auto* frame = new QFrame();
    frame->setObjectName(config.name);
    frame->setFrameShape(QFrame::StyledPanel);
    frame->setFrameShadow(QFrame::Raised);

    // Build and set layout if valid
    QLayout* layout = buildLayoutStructure(config.layout);
    if (layout) {
        frame->setLayout(layout);
    }

    return frame;
}

QWidget* SHantillyBuilder::buildTabWidget(const Sbx::Models::TabWidgetConfig& config) {
    auto* tabWidget = new QTabWidget();
    tabWidget->setObjectName(config.name);

    for (const auto& pageConfig : config.pages) {
        auto* page = new QWidget();
        page->setObjectName(pageConfig.name);

        QLayout* layout = buildLayoutStructure(pageConfig.layout);
        if (layout) {
            page->setLayout(layout);
        }

        tabWidget->addTab(page, pageConfig.title);
    }

    return tabWidget;
}

QLayout* SHantillyBuilder::buildLayoutStructure(const Sbx::Models::SbxLayoutConfig& config) {
    QLayout* layout = nullptr;
    switch (config.type) {
        case Sbx::Models::SbxLayoutConfig::VBox:
            layout = new QVBoxLayout();
            break;
        case Sbx::Models::SbxLayoutConfig::HBox:
            layout = new QHBoxLayout();
            break;
        default:
            layout = new QVBoxLayout();
            break;
    }

    if (layout) {
        layout->setSpacing(config.spacing);
        layout->setContentsMargins(config.margin, config.margin, config.margin, config.margin);
    }

    return layout;
}

PushButtonWidget* SHantillyBuilder::buildPushButton(const QString& title, const QString& name) {
    // Implementation for legacy support using the new config-based methods
    Sbx::Models::ButtonConfig config;
    config.text = title;
    config.name = name;

    // Note: buildButton currently returns QPushButton, but legacy expects PushButtonWidget.
    // This highlights that we might need to migrate the implementation details later.
    // For now, this stub exists to satisfy the interface.
    return nullptr;
}
