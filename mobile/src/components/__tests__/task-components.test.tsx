import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react-native';

import { TaskCard } from '../task-card';
import { TaskPagination } from '../task-pagination';
import { TaskSearchBar } from '../task-search-bar';
import { TaskStatusFilter } from '../task-status-filter';
import { Task } from '@/services/api';

describe('Task Management UI Component Tests', () => {
  describe('TaskSearchBar', () => {
    it('renders with placeholder and triggers onChangeText', async () => {
      const onChangeTextMock = jest.fn();
      const onClearMock = jest.fn();

      await render(
        <TaskSearchBar
          value=""
          onChangeText={onChangeTextMock}
          onClear={onClearMock}
          placeholder="Search test..."
        />
      );

      const input = screen.getByTestId('task-search-input');
      expect(input).toBeTruthy();
      expect(screen.getByPlaceholderText('Search test...')).toBeTruthy();

      await fireEvent.changeText(input, 'golang');
      expect(onChangeTextMock).toHaveBeenCalledWith('golang');
    });

    it('renders clear button when value is present and triggers onClear', async () => {
      const onChangeTextMock = jest.fn();
      const onClearMock = jest.fn();

      await render(
        <TaskSearchBar
          value="active keyword"
          onChangeText={onChangeTextMock}
          onClear={onClearMock}
        />
      );

      const clearBtn = screen.getByTestId('task-search-clear-button');
      expect(clearBtn).toBeTruthy();

      await fireEvent.press(clearBtn);
      expect(onClearMock).toHaveBeenCalledTimes(1);
    });
  });

  describe('TaskStatusFilter', () => {
    it('renders all status options and calls onSelectStatus on selection', async () => {
      const onSelectStatusMock = jest.fn();

      await render(
        <TaskStatusFilter
          selectedStatus="all"
          onSelectStatus={onSelectStatusMock}
        />
      );

      expect(screen.getByText('All')).toBeTruthy();
      expect(screen.getByText('To Do')).toBeTruthy();
      expect(screen.getByText('In Progress')).toBeTruthy();
      expect(screen.getByText('Done')).toBeTruthy();

      const inProgressPill = screen.getByTestId('status-filter-in_progress');
      await fireEvent.press(inProgressPill);

      expect(onSelectStatusMock).toHaveBeenCalledWith('in_progress');
    });
  });

  describe('TaskCard', () => {
    const sampleTask: Task = {
      id: 42,
      title: 'Implement Redis Caching',
      description: 'Add 60s TTL and cache invalidation',
      status: 'in_progress',
      assignee: 'backend_dev',
      created_at: '2026-11-19T00:00:00Z',
      updated_at: '2026-11-19T00:00:00Z',
    };

    it('displays task details, assignee, status badge and triggers edit & delete', async () => {
      const onEditMock = jest.fn();
      const onDeleteMock = jest.fn();

      await render(
        <TaskCard
          task={sampleTask}
          onEdit={onEditMock}
          onDelete={onDeleteMock}
        />
      );

      expect(screen.getByText('Implement Redis Caching')).toBeTruthy();
      expect(screen.getByText('Add 60s TTL and cache invalidation')).toBeTruthy();
      expect(screen.getByText('IN PROGRESS')).toBeTruthy();
      expect(screen.getByText('@backend_dev')).toBeTruthy();

      // Trigger Edit
      const editButton = screen.getByTestId('task-edit-button-42');
      await fireEvent.press(editButton);
      expect(onEditMock).toHaveBeenCalledWith(sampleTask);

      // Trigger Delete
      const deleteButton = screen.getByTestId('task-delete-button-42');
      await fireEvent.press(deleteButton);
      expect(onDeleteMock).toHaveBeenCalledWith(sampleTask);
    });
  });

  describe('TaskPagination', () => {
    it('renders correct page metadata and handles page navigation', async () => {
      const onPageChangeMock = jest.fn();

      await render(
        <TaskPagination
          page={1}
          totalPages={3}
          total={25}
          onPageChange={onPageChangeMock}
        />
      );

      expect(screen.getByText('Page 1 of 3')).toBeTruthy();
      expect(screen.getByText('(25 tasks)')).toBeTruthy();

      // Previous button should be disabled on page 1
      const prevBtn = screen.getByTestId('pagination-prev-button');
      await fireEvent.press(prevBtn);
      expect(onPageChangeMock).not.toHaveBeenCalled();

      // Next button should navigate to page 2
      const nextBtn = screen.getByTestId('pagination-next-button');
      await fireEvent.press(nextBtn);
      expect(onPageChangeMock).toHaveBeenCalledWith(2);
    });
  });
});
