import React, { useMemo } from 'react';
import { View, Text, StyleSheet, Dimensions } from 'react-native';
import Svg, { Rect, Text as SvgText, Line } from 'react-native-svg';
import { COLORS, SPACING, FONT_SIZES } from '@/constants/theme';

interface BarDataPoint {
  label: string;
  value: number;
  color?: string;
}

interface BarChartProps {
  data: BarDataPoint[];
  width?: number;
  height?: number;
  color?: string;
  showGrid?: boolean;
  showValues?: boolean;
  yAxisLabel?: string;
}

const CHART_WIDTH = Dimensions.get('window').width - SPACING.lg * 2;
const CHART_HEIGHT = 200;
const PADDING = { top: 20, right: 20, bottom: 50, left: 50 };

export const BarChart: React.FC<BarChartProps> = ({
  data,
  width = CHART_WIDTH,
  height = CHART_HEIGHT,
  color = COLORS.primary,
  showGrid = true,
  showValues = true,
  yAxisLabel = '',
}) => {
  const chartData = useMemo(() => {
    if (data.length === 0) return null;

    const chartWidth = width - PADDING.left - PADDING.right;
    const chartHeight = height - PADDING.top - PADDING.bottom;

    const maxValue = Math.max(...data.map(d => d.value));
    const yPadding = maxValue * 0.1;
    const scaledMax = maxValue + yPadding;

    const barWidth = chartWidth / data.length * 0.7;
    const barSpacing = chartWidth / data.length * 0.3;

    const bars = data.map((item, index) => {
      const barHeight = (item.value / scaledMax) * chartHeight;
      const x = PADDING.left + index * (barWidth + barSpacing) + barSpacing / 2;
      const y = height - PADDING.bottom - barHeight;

      return {
        x,
        y,
        width: barWidth,
        height: barHeight,
        value: item.value,
        label: item.label,
        color: item.color || color,
      };
    });

    const gridLines = Array.from({ length: 5 }, (_, i) => {
      const y = PADDING.top + (i / 4) * chartHeight;
      const value = scaledMax - (i / 4) * scaledMax;
      return { y, value };
    });

    return { bars, gridLines, maxValue: scaledMax };
  }, [data, width, height, color]);

  if (!chartData || data.length === 0) {
    return (
      <View style={[styles.container, { width, height }]}>
        <Text style={styles.emptyText}>No data available</Text>
      </View>
    );
  }

  return (
    <View style={[styles.container, { width, height }]}>
      <Svg width={width} height={height}>
        {showGrid && chartData.gridLines.map((line, index) => (
          <React.Fragment key={index}>
            <Line
              x1={PADDING.left}
              y1={line.y}
              x2={width - PADDING.right}
              y2={line.y}
              stroke={COLORS.border.light}
              strokeWidth="1"
              strokeDasharray="4,4"
            />
            <SvgText
              x={PADDING.left - 8}
              y={line.y + 4}
              fontSize={FONT_SIZES.xs}
              fill={COLORS.text.tertiary}
              textAnchor="end"
            >
              {line.value.toFixed(0)}
            </SvgText>
          </React.Fragment>
        ))}

        {chartData.bars.map((bar, index) => (
          <React.Fragment key={index}>
            <Rect
              x={bar.x}
              y={bar.y}
              width={bar.width}
              height={bar.height}
              fill={bar.color}
              rx="6"
              ry="6"
            />
            {showValues && bar.height > 20 && (
              <SvgText
                x={bar.x + bar.width / 2}
                y={bar.y - 8}
                fontSize={FONT_SIZES.xs}
                fill={COLORS.text.secondary}
                textAnchor="middle"
                fontWeight="600"
              >
                {bar.value.toFixed(0)}
              </SvgText>
            )}
            <SvgText
              x={bar.x + bar.width / 2}
              y={height - PADDING.bottom + 20}
              fontSize={FONT_SIZES.xs}
              fill={COLORS.text.tertiary}
              textAnchor="middle"
            >
              {bar.label}
            </SvgText>
          </React.Fragment>
        ))}
      </Svg>

      {yAxisLabel && (
        <View style={styles.yAxisLabelContainer}>
          <Text style={styles.yAxisLabel}>{yAxisLabel}</Text>
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    backgroundColor: COLORS.background.card,
    borderRadius: 12,
    padding: SPACING.sm,
  },
  emptyText: {
    fontSize: FONT_SIZES.sm,
    color: COLORS.text.tertiary,
    textAlign: 'center',
    marginTop: SPACING.xl,
  },
  yAxisLabelContainer: {
    position: 'absolute',
    left: 0,
    top: '50%',
    transform: [{ rotate: '-90deg' }],
  },
  yAxisLabel: {
    fontSize: FONT_SIZES.xs,
    color: COLORS.text.secondary,
    fontWeight: '600',
  },
});
