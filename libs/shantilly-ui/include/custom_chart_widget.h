#ifndef CUSTOM_CHART_WIDGET_H
#define CUSTOM_CHART_WIDGET_H

#include <QGraphicsSimpleTextItem>
#include <QMouseEvent>
#include <QtCharts/QChart>
#include <QtCharts/QChartView>
#include <QtCharts/QLogValueAxis>
#include <QtCharts/QPieSlice>
#include <QtCharts/QValueAxis>

class CustomChartWidget : public QChartView {
    Q_OBJECT

public:
    explicit CustomChartWidget(QWidget* parent = nullptr);
    ~CustomChartWidget() override;

    void setChartTitle(const QString& title);
    void addPoint(const QString& label, double value);
    void clearSeries();

    void setData(const QString& data);
    void appendData(const QString& data);
    void loadFromFile(const QString& filePath);
    void setAxis(const QString& config);
    void exportChart(const QString& path);

signals:
    void itemClicked(const QString& label);

protected:
    void mousePressEvent(QMouseEvent* event) override;

private slots:
    void onPieSeriesHovered(QPieSlice* slice, bool state);
    void handleMarkerClicked();
    void onSeriesHovered(const QPointF& point, bool state);

private:
    QChart* m_chart;
    QGraphicsSimpleTextItem* m_tooltip;
};

#endif // CUSTOM_CHART_WIDGET_H