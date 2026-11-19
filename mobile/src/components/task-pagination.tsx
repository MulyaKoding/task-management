import React from 'react';
import { StyleSheet, TouchableOpacity, View } from 'react-native';
import { ThemedText } from './themed-text';

interface TaskPaginationProps {
  page: number;
  totalPages: number;
  total: number;
  onPageChange: (newPage: number) => void;
  isLoading?: boolean;
}

export function TaskPagination({
  page,
  totalPages,
  total,
  onPageChange,
  isLoading = false,
}: TaskPaginationProps) {
  const canGoPrev = page > 1 && !isLoading;
  const canGoNext = page < totalPages && !isLoading;

  return (
    <View style={styles.container}>
      <TouchableOpacity
        testID="pagination-prev-button"
        style={[styles.button, !canGoPrev && styles.buttonDisabled]}
        onPress={() => canGoPrev && onPageChange(page - 1)}
        disabled={!canGoPrev}
        activeOpacity={0.7}>
        <ThemedText style={[styles.buttonText, !canGoPrev && styles.buttonTextDisabled]}>
          ‹ Previous
        </ThemedText>
      </TouchableOpacity>

      <View style={styles.infoContainer}>
        <ThemedText style={styles.pageText}>
          Page {totalPages > 0 ? page : 0} of {totalPages}
        </ThemedText>
        <ThemedText style={styles.totalText}>({total} tasks)</ThemedText>
      </View>

      <TouchableOpacity
        testID="pagination-next-button"
        style={[styles.button, !canGoNext && styles.buttonDisabled]}
        onPress={() => canGoNext && onPageChange(page + 1)}
        disabled={!canGoNext}
        activeOpacity={0.7}>
        <ThemedText style={[styles.buttonText, !canGoNext && styles.buttonTextDisabled]}>
          Next ›
        </ThemedText>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 14,
    paddingHorizontal: 8,
    borderTopWidth: 1,
    borderTopColor: '#EBEBEF',
    marginTop: 8,
  },
  button: {
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: 8,
    backgroundColor: '#F0F0F3',
  },
  buttonDisabled: {
    opacity: 0.4,
  },
  buttonText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#007AFF',
  },
  buttonTextDisabled: {
    color: '#8E8E93',
  },
  infoContainer: {
    alignItems: 'center',
  },
  pageText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#1C1C1E',
  },
  totalText: {
    fontSize: 11,
    color: '#8E8E93',
  },
});
