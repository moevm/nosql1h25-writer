import React, { useEffect, useState } from 'react';
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  useReactTable
} from '@tanstack/react-table';
import { Input, Pagination, Typography } from 'antd';
import type {Column, ColumnDef, FilterFn} from '@tanstack/react-table';

import './AdminUsers.css';

const { Title } = Typography;
const { Search } = Input;

interface User {
  id: string;
  displayName: string;
  email: string;
  systemRole: string;
  balance: number;
  createdAt: string;
}

const simpleStringFilter: FilterFn<User> = (row, columnId, value) => {
  const valueString = String(row.getValue(columnId) ?? '');
  return valueString.toLowerCase().includes(String(value).toLowerCase());
};

function DebouncedInput({
  value: initialValue,
  onChange,
  debounce = 500,
  ...props
}: {
  value: string | number;
  onChange: (value: string | number) => void;
  debounce?: number;
} & Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'>) {
  const [value, setValue] = useState(initialValue);

  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  useEffect(() => {
    const timeout = setTimeout(() => {
      onChange(value);
    }, debounce);

    return () => clearTimeout(timeout);
  }, [value]);

  return (
    <Input
      {...props}
      value={value}
      onChange={(e) => setValue(e.target.value)}
      size="small"
    />
  );
}

function Filter({ column }: { column: Column<User, any> }) {
  const columnFilterValue = column.getFilterValue();

  return (
    <div>
      <DebouncedInput
        value={(columnFilterValue ?? '') as string}
        onChange={(value) => column.setFilterValue(value)}
        placeholder={`Поиск...`}
        style={{ width: '100%' }}
      />
    </div>
  );
}

const columnHelper = createColumnHelper<User>();

const columns: Array<ColumnDef<User, any>> = [
  columnHelper.accessor('id', {
    header: 'ID',
    cell: (info) => info.getValue(),
    filterFn: simpleStringFilter,
  }),
  columnHelper.accessor('displayName', {
    header: 'Имя',
    cell: (info) => info.getValue(),
    filterFn: simpleStringFilter,
  }),
  columnHelper.accessor('email', {
    header: 'Email',
    cell: (info) => info.getValue(),
    filterFn: simpleStringFilter,
  }),
  columnHelper.accessor('systemRole', {
    header: 'Роль',
    cell: (info) => info.getValue(),
    filterFn: simpleStringFilter,
  }),
  columnHelper.accessor('balance', {
    header: 'Баланс',
    cell: (info) => info.getValue(),
    filterFn: simpleStringFilter,
  }),
  columnHelper.accessor('createdAt', {
    header: 'Дата регистрации',
    cell: (info) => new Date(info.getValue()).toLocaleDateString(),
    filterFn: simpleStringFilter,
  }),
];

const AdminUsers = () => {
  const [users, setUsers] = useState<Array<User>>([]);
  const [loading, setLoading] = useState(true);
  const [globalFilter, setGlobalFilter] = useState('');
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  });

  useEffect(() => {
    // TODO: Implement fetching users from API
    // Пример фейковых данных:
    const fakeUsers: Array<User> = [
      { id: '1', displayName: 'User 1', email: 'user1@example.com', systemRole: 'user', balance: 100, createdAt: new Date().toISOString() },
      { id: '2', displayName: 'Admin User', email: 'admin@example.com', systemRole: 'admin', balance: 500, createdAt: new Date().toISOString() },
      { id: '3', displayName: 'User 2', email: 'user2@example.com', systemRole: 'user', balance: 200, createdAt: new Date().toISOString() },
       { id: '4', displayName: 'Another User', email: 'another@example.com', systemRole: 'user', balance: 150, createdAt: new Date().toISOString() },
      { id: '5', displayName: 'Test Admin', email: 'testadmin@example.com', systemRole: 'admin', balance: 600, createdAt: new Date().toISOString() },
    ];
    setUsers(fakeUsers);
    setLoading(false);
  }, []);

  const table = useReactTable<User>({
    data: users,
    columns,
    state: {
      globalFilter,
      pagination,
    },
    onGlobalFilterChange: setGlobalFilter,
    onPaginationChange: setPagination,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    filterFns: {
      simpleString: simpleStringFilter,
    } as Record<string, FilterFn<User>>,
    // Добавьте getSortedRowModel и sortingFns если нужна сортировка
    // debugTable: true,
  });

  return (
    <div className="admin-users">
      <Title level={2}>Управление пользователями</Title>
      <Search
        placeholder="Глобальный поиск..."
        allowClear
        enterButton="Поиск"
        size="large"
        value={globalFilter}
        onChange={(e) => setGlobalFilter(e.target.value)}
        style={{ marginBottom: 20, width: 300 }}
      />
      <div className="table-container">
        <table>
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <th key={header.id} colSpan={header.colSpan}>
                    {header.isPlaceholder
                      ? null
                      : (
                          <>
                            <div>
                              {flexRender(
                                header.column.columnDef.header,
                                header.getContext()
                              )}
                            </div>
                            {header.column.getCanFilter() ? (
                              <div style={{ marginTop: 8 }}>
                                <Filter column={header.column} />
                              </div>
                            ) : null}
                          </>
                        )}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan={columns.length} style={{ textAlign: 'center' }}>
                  Загрузка...
                </td>
              </tr>
            ) : table.getRowModel().rows.length === 0 ? (
              <tr>
                <td colSpan={columns.length} style={{ textAlign: 'center' }}>
                  Пользователи не найдены.
                </td>
              </tr>
            ) : (
              table.getRowModel().rows.map((row) => (
                <tr key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
      <div style={{ marginTop: 20, display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
        <Pagination
          total={table.getFilteredRowModel().rows.length}
          current={table.getState().pagination.pageIndex + 1}
          pageSize={table.getState().pagination.pageSize}
          onChange={(page, pageSize) => table.setPagination({ pageIndex: page - 1, pageSize: pageSize || 10 })}
          showSizeChanger
          showTotal={(total) => `Всего ${total} пользователей`}
          pageSizeOptions={['5', '10', '20', '50']}
        />
      </div>
    </div>
  );
};

export default AdminUsers; 