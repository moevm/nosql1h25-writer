import { useEffect, useState } from 'react';
import { Card, Col, DatePicker, Input, InputNumber, Pagination, Row, Select, Space } from 'antd';
import dayjs from 'dayjs';
import type { RangePickerProps } from 'antd/es/date-picker';
import { api } from '@/integrations/api.ts';
import './users.css';

interface User {
    id: string;
    displayName: string;
    email: string;
    systemRole: 'admin' | 'user';
    balance: number;
    clientRating: number;
    freelancerRating: number;
    createdAt: string;
    updatedAt: string;
}

interface UsersResponse {
    users: Array<User>;
    total: number;
}

type SortOption = 'newest' | 'oldest' | 'rich' | 'poor' | 'name_asc' | 'name_desc' | 
                 'freelancer_rating_asc' | 'freelancer_rating_desc' | 
                 'client_rating_asc' | 'client_rating_desc';

interface Filters {
    nameSearch?: string;
    emailSearch?: string;
    roles: Array<string>;
    minFreelancerRating?: number;
    maxFreelancerRating?: number;
    minClientRating?: number;
    maxClientRating?: number;
    minCreatedAt?: string;
    maxCreatedAt?: string;
    minBalance?: number;
    maxBalance?: number;
    sortBy: SortOption;
}

const DEBOUNCE_DELAY = 500; // 500ms задержка

export const UsersList = () => {
    const [users, setUsers] = useState<Array<User>>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [filters, setFilters] = useState<Filters>({
        roles: [],
        sortBy: 'newest'
    });
    const [debouncedFilters, setDebouncedFilters] = useState<Filters>(filters);
    const [tableLoading, setTableLoading] = useState(false);

    const fetchUsers = async (page: number, size: number, filterParams: Filters) => {
        try {
            setTableLoading(true);
            const params: Record<string, any> = {
                offset: (page - 1) * size,
                limit: size,
                nameSearch: filterParams.nameSearch || undefined,
                emailSearch: filterParams.emailSearch || undefined,
                minBalance: filterParams.minBalance || undefined,
                maxBalance: filterParams.maxBalance || undefined,
                minFreelancerRating: filterParams.minFreelancerRating || undefined,
                maxFreelancerRating: filterParams.maxFreelancerRating || undefined,
                minClientRating: filterParams.minClientRating || undefined,
                maxClientRating: filterParams.maxClientRating || undefined,
                minCreatedAt: filterParams.minCreatedAt || undefined,
                maxCreatedAt: filterParams.maxCreatedAt || undefined,
                sortBy: filterParams.sortBy,
            };

            // Добавляем каждый role как отдельный параметр
            if (filterParams.roles.length > 0) {
                filterParams.roles.forEach(role => {
                    params['role'] = role;
                });
            }

            const response = await api.get<UsersResponse>('/admin/users', { params });
            setUsers(response.data.users);
            setTotal(response.data.total);
            setError(null);
        } catch (err) {
            setError('Ошибка при загрузке пользователей');
            console.error('Error fetching users:', err);
        } finally {
            setTableLoading(false);
            setLoading(false);
        }
    };

    // Эффект для начальной загрузки
    useEffect(() => {
        fetchUsers(currentPage, pageSize, filters);
    }, []);

    // Эффект для debounce фильтров
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedFilters(filters);
        }, DEBOUNCE_DELAY);

        return () => clearTimeout(timer);
    }, [filters]);

    // Эффект для загрузки данных при изменении debounced фильтров
    useEffect(() => {
        if (!loading) {
            fetchUsers(currentPage, pageSize, debouncedFilters);
        }
    }, [currentPage, pageSize, debouncedFilters, loading]);

    const handlePageChange = (page: number, size: number) => {
        setCurrentPage(page);
        setPageSize(size);
    };

    const handleFilterChange = (key: keyof Filters, value: any) => {
        const newFilters = { ...filters, [key]: value };
        setFilters(newFilters);
        setCurrentPage(1);
    };

    const disabledDate: RangePickerProps['disabledDate'] = (current) => {
        return current > dayjs().endOf('day');
    };

    if (loading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div className="error">{error}</div>;
    }

    return (
        <div className="users-list">
            <h2>Список пользователей ({total})</h2>
            <Card className="filters-container">
                <Row gutter={[16, 16]}>
                    <Col span={24}>
                        <Space size="middle">
                            <Input.Search
                                placeholder="Поиск по имени..."
                                allowClear
                                value={filters.nameSearch}
                                onChange={(e) => handleFilterChange('nameSearch', e.target.value)}
                                onSearch={(value) => handleFilterChange('nameSearch', value)}
                                style={{ width: 200 }}
                            />
                            <Input.Search
                                placeholder="Поиск по email..."
                                allowClear
                                value={filters.emailSearch}
                                onChange={(e) => handleFilterChange('emailSearch', e.target.value)}
                                onSearch={(value) => handleFilterChange('emailSearch', value)}
                                style={{ width: 200 }}
                            />
                            <Select
                                mode="multiple"
                                placeholder="Роли"
                                value={filters.roles}
                                onChange={(value) => handleFilterChange('roles', value)}
                                style={{ width: 200 }}
                                options={[
                                    { value: 'user', label: 'Пользователь' },
                                    { value: 'admin', label: 'Администратор' },
                                ]}
                            />
                            <Select
                                value={filters.sortBy}
                                onChange={(value) => handleFilterChange('sortBy', value)}
                                style={{ width: 200 }}
                                options={[
                                    { value: 'newest', label: 'Сначала новые' },
                                    { value: 'oldest', label: 'Сначала старые' },
                                    { value: 'rich', label: 'Сначала богатые' },
                                    { value: 'poor', label: 'Сначала бедные' },
                                    { value: 'name_asc', label: 'Имя по возрастанию' },
                                    { value: 'name_desc', label: 'Имя по убыванию' },
                                    { value: 'freelancer_rating_asc', label: 'Рейтинг фрилансера по возрастанию' },
                                    { value: 'freelancer_rating_desc', label: 'Рейтинг фрилансера по убыванию' },
                                    { value: 'client_rating_asc', label: 'Рейтинг клиента по возрастанию' },
                                    { value: 'client_rating_desc', label: 'Рейтинг клиента по убыванию' },
                                ]}
                            />
                        </Space>
                    </Col>
                    <Col span={24}>
                        <Space size="middle">
                            <div className="filter-group">
                                <div>Рейтинг фрилансера:</div>
                                <Space>
                                    <InputNumber
                                        min={0}
                                        max={5}
                                        step={0.1}
                                        placeholder="Мин"
                                        value={filters.minFreelancerRating}
                                        onChange={(value) => handleFilterChange('minFreelancerRating', value)}
                                    />
                                    <InputNumber
                                        min={0}
                                        max={5}
                                        step={0.1}
                                        placeholder="Макс"
                                        value={filters.maxFreelancerRating}
                                        onChange={(value) => handleFilterChange('maxFreelancerRating', value)}
                                    />
                                </Space>
                            </div>
                            <div className="filter-group">
                                <div>Рейтинг клиента:</div>
                                <Space>
                                    <InputNumber
                                        min={0}
                                        max={5}
                                        step={0.1}
                                        placeholder="Мин"
                                        value={filters.minClientRating}
                                        onChange={(value) => handleFilterChange('minClientRating', value)}
                                    />
                                    <InputNumber
                                        min={0}
                                        max={5}
                                        step={0.1}
                                        placeholder="Макс"
                                        value={filters.maxClientRating}
                                        onChange={(value) => handleFilterChange('maxClientRating', value)}
                                    />
                                </Space>
                            </div>
                            <div className="filter-group">
                                <div>Баланс:</div>
                                <Space>
                                    <InputNumber
                                        min={0}
                                        placeholder="Мин"
                                        value={filters.minBalance}
                                        onChange={(value) => handleFilterChange('minBalance', value)}
                                    />
                                    <InputNumber
                                        min={0}
                                        placeholder="Макс"
                                        value={filters.maxBalance}
                                        onChange={(value) => handleFilterChange('maxBalance', value)}
                                    />
                                </Space>
                            </div>
                            <div className="filter-group">
                                <div>Дата регистрации:</div>
                                <Space>
                                    <DatePicker
                                        placeholder="От"
                                        value={filters.minCreatedAt ? dayjs(filters.minCreatedAt) : null}
                                        onChange={(date) => handleFilterChange('minCreatedAt', date.toISOString())}
                                        disabledDate={disabledDate}
                                    />
                                    <DatePicker
                                        placeholder="До"
                                        value={filters.maxCreatedAt ? dayjs(filters.maxCreatedAt) : null}
                                        onChange={(date) => handleFilterChange('maxCreatedAt', date.toISOString())}
                                        disabledDate={disabledDate}
                                    />
                                </Space>
                            </div>
                        </Space>
                    </Col>
                </Row>
            </Card>
            <div className="table-container">
                <table>
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Имя</th>
                            <th>Email</th>
                            <th>Роль</th>
                            <th>Баланс</th>
                            <th>Рейтинг клиента</th>
                            <th>Рейтинг фрилансера</th>
                            <th>Дата создания</th>
                        </tr>
                    </thead>
                    <tbody>
                        {users.map((user) => (
                            <tr key={user.id}>
                                <td>{user.id}</td>
                                <td>{user.displayName}</td>
                                <td>{user.email}</td>
                                <td>{user.systemRole}</td>
                                <td>{user.balance}</td>
                                <td>{user.clientRating.toFixed(1)}</td>
                                <td>{user.freelancerRating.toFixed(1)}</td>
                                <td>{new Date(user.createdAt).toLocaleDateString()}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {tableLoading && <div className="table-loading">Загрузка...</div>}
            </div>
            <div className="pagination-container">
                <Pagination
                    current={currentPage}
                    pageSize={pageSize}
                    total={total}
                    onChange={handlePageChange}
                    showSizeChanger
                    showTotal={(totalCount) => `Всего ${totalCount} пользователей`}
                    pageSizeOptions={['5', '10', '20', '50', '100']}
                />
            </div>
        </div>
    );
}; 