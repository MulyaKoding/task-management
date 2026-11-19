import React, { useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Modal,
  Platform,
  ScrollView,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { ApiError, Task, TaskApi, TaskStatus } from '@/services/api';
import { ThemedText } from './themed-text';

interface TaskModalProps {
  visible: boolean;
  task?: Task | null;
  onClose: () => void;
  onSaved: () => void;
}

const STATUSES: { label: string; value: TaskStatus }[] = [
  { label: 'To Do', value: 'todo' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Done', value: 'done' },
];

export function TaskModal({ visible, task, onClose, onSaved }: TaskModalProps) {
  const isEditing = Boolean(task);

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState<TaskStatus>('todo');
  const [assignee, setAssignee] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    if (task) {
      setTitle(task.title || '');
      setDescription(task.description || '');
      setStatus(task.status || 'todo');
      setAssignee(task.assignee || '');
    } else {
      setTitle('');
      setDescription('');
      setStatus('todo');
      setAssignee('');
    }
    setErrorMessage(null);
    setIsSubmitting(false);
  }, [task, visible]);

  const handleSubmit = async () => {
    if (!title.trim()) {
      setErrorMessage('Task title is required.');
      return;
    }

    setIsSubmitting(true);
    setErrorMessage(null);

    try {
      if (isEditing && task) {
        await TaskApi.updateTask(task.id, {
          title: title.trim(),
          description: description.trim() || undefined,
          status,
          assignee: assignee.trim() || undefined,
        });
      } else {
        await TaskApi.createTask({
          title: title.trim(),
          description: description.trim() || undefined,
          status,
          assignee: assignee.trim() || undefined,
        });
      }

      onSaved();
      onClose();
    } catch (err: any) {
      if (err instanceof ApiError) {
        if (err.status === 409 || err.code === 'CONFLICT') {
          setErrorMessage('Conflict (409): A task with this title already exists. Please choose a different title.');
        } else {
          setErrorMessage(err.message || 'Operation failed. Please try again.');
        }
      } else {
        setErrorMessage(err?.message || 'An unexpected error occurred.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Modal
      testID="task-modal"
      visible={visible}
      animationType="slide"
      transparent={true}
      onRequestClose={onClose}>
      <View style={styles.overlay}>
        <View style={styles.dialogContainer}>
          <ScrollView contentContainerStyle={styles.scrollContent}>
            <View style={styles.header}>
              <ThemedText style={styles.headerTitle}>
                {isEditing ? 'Edit Task' : 'Create New Task'}
              </ThemedText>
              <TouchableOpacity
                testID="task-modal-close-button"
                onPress={onClose}
                disabled={isSubmitting}
                style={styles.closeBtn}>
                <ThemedText style={styles.closeBtnText}>✕</ThemedText>
              </TouchableOpacity>
            </View>

            {errorMessage ? (
              <View style={styles.errorBanner} testID="task-modal-error">
                <ThemedText style={styles.errorText}>{errorMessage}</ThemedText>
              </View>
            ) : null}

            <View style={styles.formGroup}>
              <ThemedText style={styles.label}>
                Title <ThemedText style={styles.required}>*</ThemedText>
              </ThemedText>
              <TextInput
                testID="task-modal-title-input"
                style={styles.input}
                placeholder="e.g. Implement user authentication"
                placeholderTextColor="#8E8E93"
                value={title}
                onChangeText={setTitle}
                editable={!isSubmitting}
              />
            </View>

            <View style={styles.formGroup}>
              <ThemedText style={styles.label}>Description</ThemedText>
              <TextInput
                testID="task-modal-description-input"
                style={[styles.input, styles.textArea]}
                placeholder="Provide task details or acceptance criteria..."
                placeholderTextColor="#8E8E93"
                value={description}
                onChangeText={setDescription}
                multiline={true}
                numberOfLines={3}
                editable={!isSubmitting}
              />
            </View>

            <View style={styles.formGroup}>
              <ThemedText style={styles.label}>Status</ThemedText>
              <View style={styles.statusRow}>
                {STATUSES.map((st) => {
                  const isSelected = status === st.value;
                  return (
                    <TouchableOpacity
                      key={st.value}
                      testID={`task-modal-status-${st.value}`}
                      style={[styles.statusOption, isSelected && styles.statusOptionSelected]}
                      onPress={() => setStatus(st.value)}
                      disabled={isSubmitting}
                      activeOpacity={0.8}>
                      <ThemedText
                        style={[
                          styles.statusOptionText,
                          isSelected && styles.statusOptionTextSelected,
                        ]}>
                        {st.label}
                      </ThemedText>
                    </TouchableOpacity>
                  );
                })}
              </View>
            </View>

            <View style={styles.formGroup}>
              <ThemedText style={styles.label}>Assignee</ThemedText>
              <TextInput
                testID="task-modal-assignee-input"
                style={styles.input}
                placeholder="e.g. john_doe"
                placeholderTextColor="#8E8E93"
                value={assignee}
                onChangeText={setAssignee}
                autoCapitalize="none"
                editable={!isSubmitting}
              />
            </View>

            <View style={styles.buttonRow}>
              <TouchableOpacity
                testID="task-modal-cancel-button"
                style={styles.cancelButton}
                onPress={onClose}
                disabled={isSubmitting}
                activeOpacity={0.7}>
                <ThemedText style={styles.cancelButtonText}>Cancel</ThemedText>
              </TouchableOpacity>

              <TouchableOpacity
                testID="task-modal-save-button"
                style={[styles.saveButton, isSubmitting && styles.saveButtonDisabled]}
                onPress={handleSubmit}
                disabled={isSubmitting}
                activeOpacity={0.8}>
                {isSubmitting ? (
                  <ActivityIndicator color="#FFFFFF" size="small" />
                ) : (
                  <ThemedText style={styles.saveButtonText}>
                    {isEditing ? 'Save Changes' : 'Create Task'}
                  </ThemedText>
                )}
              </TouchableOpacity>
            </View>
          </ScrollView>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    justifyContent: 'center',
    alignItems: 'center',
    padding: 16,
  },
  dialogContainer: {
    width: '100%',
    maxWidth: 500,
    backgroundColor: '#FFFFFF',
    borderRadius: 16,
    overflow: 'hidden',
    maxHeight: '90%',
    ...Platform.select({
      ios: {
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.2,
        shadowRadius: 10,
      },
      android: {
        elevation: 8,
      },
    }),
  },
  scrollContent: {
    padding: 20,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#1C1C1E',
  },
  closeBtn: {
    padding: 6,
  },
  closeBtnText: {
    fontSize: 18,
    color: '#8E8E93',
  },
  errorBanner: {
    backgroundColor: '#FFEBEA',
    borderColor: '#FF4D4F',
    borderWidth: 1,
    borderRadius: 8,
    padding: 10,
    marginBottom: 16,
  },
  errorText: {
    color: '#D9363E',
    fontSize: 13,
    fontWeight: '500',
  },
  formGroup: {
    marginBottom: 16,
  },
  label: {
    fontSize: 14,
    fontWeight: '600',
    color: '#2C3E50',
    marginBottom: 6,
  },
  required: {
    color: '#FF3B30',
  },
  input: {
    backgroundColor: '#F5F5F7',
    borderWidth: 1,
    borderColor: '#E5E5EA',
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 10,
    fontSize: 15,
    color: '#000000',
  },
  textArea: {
    height: 80,
    textAlignVertical: 'top',
  },
  statusRow: {
    flexDirection: 'row',
    gap: 8,
  },
  statusOption: {
    flex: 1,
    paddingVertical: 8,
    borderRadius: 8,
    backgroundColor: '#F0F0F3',
    alignItems: 'center',
    borderWidth: 1,
    borderColor: '#E2E2E6',
  },
  statusOptionSelected: {
    backgroundColor: '#007AFF',
    borderColor: '#007AFF',
  },
  statusOptionText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#60646C',
  },
  statusOptionTextSelected: {
    color: '#FFFFFF',
  },
  buttonRow: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 12,
    marginTop: 12,
  },
  cancelButton: {
    paddingVertical: 10,
    paddingHorizontal: 16,
    borderRadius: 8,
    backgroundColor: '#F2F2F7',
  },
  cancelButtonText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#60646C',
  },
  saveButton: {
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 8,
    backgroundColor: '#007AFF',
    minWidth: 120,
    alignItems: 'center',
    justifyContent: 'center',
  },
  saveButtonDisabled: {
    opacity: 0.6,
  },
  saveButtonText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#FFFFFF',
  },
});
