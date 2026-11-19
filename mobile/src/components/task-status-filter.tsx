import React from 'react';
import { StyleSheet, TouchableOpacity, View } from 'react-native';
import { ThemedText } from './themed-text';

export type FilterStatus = 'all' | 'todo' | 'in_progress' | 'done';

interface TaskStatusFilterProps {
  selectedStatus: FilterStatus;
  onSelectStatus: (status: FilterStatus) => void;
}

const STATUS_OPTIONS: { label: string; value: FilterStatus }[] = [
  { label: 'All', value: 'all' },
  { label: 'To Do', value: 'todo' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Done', value: 'done' },
];

export function TaskStatusFilter({ selectedStatus, onSelectStatus }: TaskStatusFilterProps) {
  return (
    <View style={styles.container}>
      {STATUS_OPTIONS.map((option) => {
        const isSelected = selectedStatus === option.value;
        return (
          <TouchableOpacity
            key={option.value}
            testID={`status-filter-${option.value}`}
            style={[styles.pill, isSelected && styles.pillSelected]}
            onPress={() => onSelectStatus(option.value)}
            activeOpacity={0.8}>
            <ThemedText style={[styles.pillText, isSelected && styles.pillTextSelected]}>
              {option.label}
            </ThemedText>
          </TouchableOpacity>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    gap: 8,
    marginBottom: 16,
    flexWrap: 'wrap',
  },
  pill: {
    paddingVertical: 6,
    paddingHorizontal: 14,
    borderRadius: 20,
    backgroundColor: '#F0F0F3',
    borderWidth: 1,
    borderColor: '#E2E2E6',
  },
  pillSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  pillText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#60646C',
  },
  pillTextSelected: {
    color: '#FFFFFF',
  },
});
