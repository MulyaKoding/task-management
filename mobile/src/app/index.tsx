import React, { useCallback, useEffect, useState } from 'react';
import {
  ActivityIndicator,
  Alert,
  FlatList,
  Platform,
  RefreshControl,
  StyleSheet,
  TouchableOpacity,
  View,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { TaskCard } from '@/components/task-card';
import { TaskModal } from '@/components/task-modal';
import { TaskPagination } from '@/components/task-pagination';
import { TaskSearchBar } from '@/components/task-search-bar';
import { FilterStatus, TaskStatusFilter } from '@/components/task-status-filter';
import { ThemedText } from '@/components/themed-text';
import { ThemedView } from '@/components/themed-view';
import { BottomTabInset, MaxContentWidth, Spacing } from '@/constants/theme';
import { PaginationMeta, Task, TaskApi } from '@/services/api';

export default function HomeScreen() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [meta, setMeta] = useState<PaginationMeta>({
    page: 1,
    limit: 10,
    total: 0,
    total_pages: 0,
  });

  // Filter & Search states
  const [keyword, setKeyword] = useState('');
  const [status, setStatus] = useState<FilterStatus>('all');
  const [page, setPage] = useState(1);

  // UI states
  const [isLoading, setIsLoading] = useState(false);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // Modal states
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | null>(null);

  // Load tasks from backend
  const loadTasks = useCallback(
    async (currentPage = page, currentStatus = status, currentKeyword = keyword) => {
      setIsLoading(true);
      setErrorMessage(null);

      try {
        const filterStatus = currentStatus === 'all' ? undefined : currentStatus;
        const res = await TaskApi.fetchTasks({
          page: currentPage,
          limit: 10,
          status: filterStatus,
          keyword: currentKeyword.trim() || undefined,
          sort: 'created_at DESC',
        });

        setTasks(res.data);
        setMeta(res.meta);
      } catch (err: any) {
        setErrorMessage(err?.message || 'Failed to connect to backend server.');
      } finally {
        setIsLoading(false);
        setIsRefreshing(false);
      }
    },
    [page, status, keyword]
  );

  useEffect(() => {
    loadTasks(page, status, keyword);
  }, [page, status, keyword]);

  const handleRefresh = () => {
    setIsRefreshing(true);
    loadTasks(page, status, keyword);
  };

  const handleSearchChange = (text: string) => {
    setKeyword(text);
    setPage(1); // reset to page 1 on new search
  };

  const handleClearSearch = () => {
    setKeyword('');
    setPage(1);
  };

  const handleStatusChange = (newStatus: FilterStatus) => {
    setStatus(newStatus);
    setPage(1); // reset to page 1 on filter change
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const handleOpenCreateModal = () => {
    setSelectedTask(null);
    setIsModalVisible(true);
  };

  const handleOpenEditModal = (task: Task) => {
    setSelectedTask(task);
    setIsModalVisible(true);
  };

  const handleDeleteTask = (task: Task) => {
    const confirmDelete = async () => {
      try {
        setIsLoading(true);
        await TaskApi.deleteTask(task.id);
        // Refresh list after soft-delete
        await loadTasks(page, status, keyword);
      } catch (err: any) {
        setErrorMessage(err?.message || 'Failed to delete task.');
      } finally {
        setIsLoading(false);
      }
    };

    if (Platform.OS === 'web') {
      if (typeof window !== 'undefined' && window.confirm(`Are you sure you want to delete task "${task.title}"?`)) {
        confirmDelete();
      }
    } else {
      Alert.alert(
        'Delete Task',
        `Are you sure you want to delete "${task.title}"?`,
        [
          { text: 'Cancel', style: 'cancel' },
          { text: 'Delete', style: 'destructive', onPress: confirmDelete },
        ]
      );
    }
  };

  // Called after create/update completes in modal -> refreshes list
  const handleTaskSaved = () => {
    loadTasks(page, status, keyword);
  };

  return (
    <ThemedView style={styles.container}>
      <SafeAreaView style={styles.safeArea}>
        <View style={styles.contentWrapper}>
          {/* Header */}
          <View style={styles.header}>
            <View>
              <ThemedText style={styles.title}>Tasks</ThemedText>
              <ThemedText style={styles.subtitle}>
                Manage team goals, issues & progress
              </ThemedText>
            </View>

            <TouchableOpacity
              testID="create-task-button"
              style={styles.createButton}
              onPress={handleOpenCreateModal}
              activeOpacity={0.8}>
              <ThemedText style={styles.createButtonText}>+ New Task</ThemedText>
            </TouchableOpacity>
          </View>

          {/* Search Input */}
          <TaskSearchBar
            value={keyword}
            onChangeText={handleSearchChange}
            onClear={handleClearSearch}
          />

          {/* Status Filter */}
          <TaskStatusFilter
            selectedStatus={status}
            onSelectStatus={handleStatusChange}
          />

          {/* Error Message Banner */}
          {errorMessage ? (
            <View style={styles.errorBox}>
              <ThemedText style={styles.errorBoxText}>{errorMessage}</ThemedText>
              <TouchableOpacity
                style={styles.retryButton}
                onPress={() => loadTasks(page, status, keyword)}>
                <ThemedText style={styles.retryButtonText}>Retry</ThemedText>
              </TouchableOpacity>
            </View>
          ) : null}

          {/* Loading Indicator */}
          {isLoading && !isRefreshing ? (
            <View style={styles.loadingContainer} testID="tasks-loading-indicator">
              <ActivityIndicator size="large" color="#007AFF" />
              <ThemedText style={styles.loadingText}>Loading tasks...</ThemedText>
            </View>
          ) : (
            /* Task List */
            <FlatList
              testID="tasks-flat-list"
              data={tasks}
              keyExtractor={(item) => item.id.toString()}
              renderItem={({ item }) => (
                <TaskCard
                  task={item}
                  onEdit={handleOpenEditModal}
                  onDelete={handleDeleteTask}
                />
              )}
              refreshControl={
                <RefreshControl refreshing={isRefreshing} onRefresh={handleRefresh} />
              }
              ListEmptyComponent={
                <View style={styles.emptyContainer} testID="tasks-empty-state">
                  <ThemedText style={styles.emptyTitle}>No tasks found</ThemedText>
                  <ThemedText style={styles.emptySubtitle}>
                    {keyword || status !== 'all'
                      ? 'Try clearing or modifying your search or filter'
                      : 'Get started by creating your first task!'}
                  </ThemedText>
                </View>
              }
              contentContainerStyle={styles.listContent}
            />
          )}

          {/* Pagination Controls */}
          <TaskPagination
            page={meta.page}
            totalPages={meta.total_pages}
            total={meta.total}
            onPageChange={handlePageChange}
            isLoading={isLoading}
          />
        </View>

        {/* Create / Edit Modal */}
        <TaskModal
          visible={isModalVisible}
          task={selectedTask}
          onClose={() => setIsModalVisible(false)}
          onSaved={handleTaskSaved}
        />
      </SafeAreaView>
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'center',
    backgroundColor: '#F8F9FA',
  },
  safeArea: {
    flex: 1,
    maxWidth: MaxContentWidth,
    paddingHorizontal: Spacing.four,
    paddingBottom: BottomTabInset + Spacing.two,
  },
  contentWrapper: {
    flex: 1,
    paddingTop: 8,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
  },
  title: {
    fontSize: 28,
    fontWeight: '800',
    color: '#1C1C1E',
  },
  subtitle: {
    fontSize: 13,
    color: '#8E8E93',
    marginTop: 2,
  },
  createButton: {
    backgroundColor: '#007AFF',
    paddingVertical: 9,
    paddingHorizontal: 16,
    borderRadius: 8,
  },
  createButtonText: {
    color: '#FFFFFF',
    fontWeight: '700',
    fontSize: 14,
  },
  errorBox: {
    backgroundColor: '#FFEBEA',
    borderColor: '#FF4D4F',
    borderWidth: 1,
    borderRadius: 8,
    padding: 12,
    marginBottom: 12,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  errorBoxText: {
    color: '#D9363E',
    fontSize: 13,
    flex: 1,
  },
  retryButton: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    backgroundColor: '#D9363E',
    borderRadius: 4,
    marginLeft: 8,
  },
  retryButtonText: {
    color: '#FFFFFF',
    fontSize: 12,
    fontWeight: '600',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: 40,
  },
  loadingText: {
    marginTop: 12,
    color: '#8E8E93',
    fontSize: 14,
  },
  listContent: {
    paddingBottom: 16,
  },
  emptyContainer: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 50,
  },
  emptyTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#2C3E50',
    marginBottom: 6,
  },
  emptySubtitle: {
    fontSize: 14,
    color: '#8E8E93',
    textAlign: 'center',
  },
});
