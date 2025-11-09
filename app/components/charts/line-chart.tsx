import React, { useMemo } from 'react';
import { View, Text, StyleSheet, Dimensions } from 'react-native';
import Svg, { Line, Circle, Path, Text as SvgText, Defs, LinearGradient, Stop } from 'react-native-svg';
import { COLORS, SPACING, FONT_SIZES } from '@/constants/theme';

interface DataPoint {
  x: number; // timestamp or index
  y: number; // value
  label?: string;
}

interface LineChartProps {
  data: DataPoint[];
  width?: number;
  height?: number;
  color?: string;
  showGrid?: boolean;
  showLabels?: boolean;
  yAxisLabel?: string;
  showGradient?: boolean;
}

const CHART_WIDTH = Dimensions.get('window').width - SPACING.lg * 2;
const CHART_HEIGHT = 200;
const PADDING = { top: 20, right: 20, bottom: 40, left: 50 };

export const LineChart: React.FC<LineChartProps> = ({
  data,
  width = CHART_WIDTH,
  height = CHART_HEIGHT,
  color = COLORS.primary,
  showGrid = true,
  showLabels = true,
  yAxisLabel = '',
  showGradient = true,
}) => {
  const chartData = useMemo(() => {
    if (data.length === 0) return null;

    // Calculate min/max values with padding
    const yValues = data.map(d => d.y);
    const minY = Math.min(...yValues);
    const maxY = Math.max(...yValues);
    const yRange = maxY - minY || 1;
    const yPadding = yRange * 0.1;

    const chartWidth = width - PADDING.left - PADDING.right;
    const chartHeight = height - PADDING.top - PADDING.bottom;

    // Scale functions
    const scaleX = (index: number) => {
      return PADDING.left + (index / Math.max(1, data.length - 1)) * chartWidth;
    };

    const scaleY = (value: number) => {
      const normalizedValue = (value - (minY - yPadding)) / (yRange + yPadding * 2);
      return height - PADDING.bottom - normalizedValue * chartHeight;
    };

    // Generate path
    const pathData = data.map((point, index) => {
      const x = scaleX(index);
      const y = scaleY(point.y);
      return index === 0 ? `M ${x} ${y}` : `L ${x} ${y}`;
    }).join(' ');

    // Generate gradient fill path
    const gradientPath = data.length > 0 
      ? `${pathData} L ${scaleX(data.length - 1)} ${height - PADDING.bottom} L ${PADDING.left} ${height - PADDING.bottom} Z`
      : '';

    // Generate grid lines (5 horizontal lines)
    const gridLines = Array.from({ length: 5 }, (_, i) => {
      const y = PADDING.top + (i / 4) * chartHeight;
      const value = maxY + yPadding - (i / 4) * (yRange + yPadding * 2);
      return { y, value };
    });

    // Generate points
    const points = data.map((point, index) => ({
      x: scaleX(index),
      y: scaleY(point.y),
      value: point.y,
      label: point.label,
    }));

    return {
      pathData,
      gradientPath,
      points,
      gridLines,
      minY: minY - yPadding,
      maxY: maxY + yPadding,
    };
  }, [data, width, height]);

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
        <Defs>
          <LinearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
            <Stop offset="0%" stopColor={color} stopOpacity="0.3" />
            <Stop offset="100%" stopColor={color} stopOpacity="0.05" />
          </LinearGradient>
        </Defs>

        {/* Grid lines */}
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
            {showLabels && (
              <SvgText
                x={PADDING.left - 8}
                y={line.y + 4}
                fontSize={FONT_SIZES.xs}
                fill={COLORS.text.tertiary}
                textAnchor="end"
              >
                {line.value.toFixed(0)}
              </SvgText>
            )}
          </React.Fragment>
        ))}

        {/* Gradient fill */}
        {showGradient && (
          <Path
            d={chartData.gradientPath}
            fill="url(#chartGradient)"
          />
        )}

        {/* Line path */}
        <Path
          d={chartData.pathData}
          stroke={color}
          strokeWidth="3"
          fill="none"
          strokeLinecap="round"
          strokeLinejoin="round"
        />

        {/* Data points */}
        {chartData.points.map((point, index) => (
          <React.Fragment key={index}>
            <Circle
              cx={point.x}
              cy={point.y}
              r="6"
              fill={COLORS.background.primary}
              stroke={color}
              strokeWidth="2"
            />
            <Circle
              cx={point.x}
              cy={point.y}
              r="3"
              fill={color}
            />
          </React.Fragment>
        ))}

        {/* X-axis labels */}
        {showLabels && chartData.points.map((point, index) => {
          // Show every other label to avoid crowding
          if (index % Math.ceil(data.length / 6) !== 0 && index !== data.length - 1) {
            return null;
          }
          return (
            <SvgText
              key={`label-${index}`}
              x={point.x}
              y={height - PADDING.bottom + 20}
              fontSize={FONT_SIZES.xs}
              fill={COLORS.text.tertiary}
              textAnchor="middle"
            >
              {point.label || index + 1}
            </SvgText>
          );
        })}
      </Svg>

      {/* Y-axis label */}
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
