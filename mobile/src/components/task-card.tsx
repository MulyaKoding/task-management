import React from 'react';
import { StyleSheet, TouchableOpacity, View } from 'react-native';
import { Task, TaskStatus } from '@/services/api';
import { ThemedText } from './themed-text';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
}

const STATUS_CONFIG: Record<TaskStatus, { label: string; bg: string; text: string }> = {
  todo: { label: 'TO DO', bg: '#FFF3D6', text: '#B7791F' },
  in_progress: { label: 'IN PROGRESS', bg: '#E1EFFF', text: '#1E6FE1' },
  done: { label: 'DONE', bg: '#E3F9E5', text: '#198038' },
};

export function TaskCard({ task, onEdit, onDelete }: TaskCardProps) {
  const statusInfo = STATUS_CONFIG[task.status] || STATUS_CONFIG.todo;

  return (
    <View style={styles.card} testID={`task-card-${task.id}`}>
      <View style={styles.headerRow}>
        <ThemedText style={styles.title} numberOfLines={2}>
          {task.title}
        </ThemedText>
        <View style={[styles.statusBadge, { backgroundColor: statusInfo.bg }]}>
          <ThemedText style={[styles.statusBadgeText, { color: statusInfo.text }]}>
            {statusInfo.label}
          </ThemedText>
        </View>
      </View>

      {task.description ? (
        <ThemedText style={styles.description} numberOfLines={3}>
          {task.description}
        </ThemedText>
      ) : null}

      <View style={styles.footerRow}>
        <View style={styles.metaCol}>
          {task.assignee ? (
            <View style={styles.assigneeContainer}>
              <ThemedText style={styles.assigneeLabel}>Assignee:</ThemedText>
              <ThemedText style={styles.assigneeValue}>@{task.assignee}</ThemedText>
            </View>
          ) : (
            <ThemedText style={styles.unassigned}>Unassigned</ThemedText>
          )}
        </View>

        <View style={styles.actions}>
          <TouchableOpacity
            testID={`task-edit-button-${task.id}`}
            style={styles.editButton}
            onPress={() => onEdit(task)}
            activeOpacity={0.7}>
            <ThemedText style={styles.editButtonText}>Edit</ThemedText>
          </TouchableOpacity>

          <TouchableOpacity
            testID={`task-delete-button-${task.id}`}
            style={styles.deleteButton}
            onPress={() => onDelete(task)}
            activeOpacity={0.7}>
            <ThemedText style={styles.deleteButtonText}>Delete</ThemedText>
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#EBEBEF',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 3,
    elevation: 2,
  },
  headerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    gap: 8,
    marginBottom: 8,
  },
  title: {
    flex: 1,
    fontSize: 16,
    fontWeight: '700',
    color: '#1C1C1E',
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 6,
  },
  statusBadgeText: {
    fontSize: 11,
    fontWeight: '700',
    letterSpacing: 0.5,
  },
  description: {
    fontSize: 14,
    color: '#60646C',
    lineHeight: 20,
    marginBottom: 12,
  },
  footerRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: 10,
    borderTopWidth: 1,
    borderTopColor: '#F2F2F7',
  },
  metaCol: {
    flex: 1,
  },
  assigneeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  assigneeLabel: {
    fontSize: 12,
    color: '#8E8E93',
  },
  assigneeValue: {
    fontSize: 12,
    fontWeight: '600',
    color: '#2C3E50',
  },
  unassigned: {
    fontSize: 12,
    fontStyle: 'italic',
    color: '#8E8E93',
  },
  actions: {
    flexDirection: 'row',
    gap: 8,
  },
  editButton: {
    paddingVertical: 5,
    paddingHorizontal: 12,
    backgroundColor: '#F0F0F3',
    borderRadius: 6,
  },
  editButtonText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#007AFF',
  },
  deleteButton: {
    paddingVertical: 5,
    paddingHorizontal: 12,
    backgroundColor: '#FFF0F0',
    borderRadius: 6,
  },
  deleteButtonText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#FF3B30',
  },
});
