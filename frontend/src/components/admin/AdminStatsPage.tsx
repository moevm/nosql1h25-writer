import { useMemo, useState } from 'react';
import axios from 'axios';
import { Alert, Button, Select, Spin } from 'antd';
import {
  CategoryScale,
  Chart as ChartJS,
  Legend as ChartLegend,
  Tooltip as ChartTooltip,
  LineElement,
  LinearScale,
  PointElement,
  TimeScale,
  Title,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import 'chartjs-adapter-date-fns';
import { differenceInDays, differenceInMonths, differenceInYears, format } from 'date-fns';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  ChartTooltip,
  ChartLegend,
  TimeScale
);

const allowedX = [
  { value: 'user_id', label: 'ID пользователя' },
  { value: 'user_system_role', label: 'Роль пользователя' },
  { value: 'user_active', label: 'Активен пользователь' },
  { value: 'user_created_at', label: 'Дата регистрации пользователя' },
  { value: 'order_id', label: 'ID заказа' },
  { value: 'order_active', label: 'Активен заказ' },
  { value: 'order_freelancer_id', label: 'ID фрилансера' },
  { value: 'order_client_id', label: 'ID клиента' },
  { value: 'order_created_at', label: 'Дата создания заказа' },
];

const allowedY = [
  { value: 'count', label: 'Количество' },
  { value: 'user_balance', label: 'Баланс пользователя' },
  { value: 'user_client_rating', label: 'Рейтинг клиента' },
  { value: 'user_freelancer_rating', label: 'Рейтинг фрилансера' },
  { value: 'order_completion_time', label: 'Время выполнения заказа' },
  { value: 'order_cost', label: 'Стоимость заказа' },
  { value: 'order_responses_count', label: 'Число откликов на заказ' },
];

type AggKey =
  | 'count'
  | 'user_balance'
  | 'user_client_rating'
  | 'user_freelancer_rating'
  | 'order_completion_time'
  | 'order_cost'
  | 'order_responses_count';

const allowedAgg: Record<AggKey, Array<{ value: string; label: string }>> = {
  count: [{ value: 'count', label: 'Количество' }],
  user_balance: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
    { value: 'sum', label: 'Сумма' },
  ],
  user_client_rating: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
  ],
  user_freelancer_rating: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
  ],
  order_completion_time: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
    { value: 'sum', label: 'Сумма' },
  ],
  order_cost: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
    { value: 'sum', label: 'Сумма' },
  ],
  order_responses_count: [
    { value: 'avg', label: 'Среднее' },
    { value: 'min', label: 'Минимум' },
    { value: 'max', label: 'Максимум' },
    { value: 'sum', label: 'Сумма' },
  ],
};

function smartGroupByDate(data: Array<{ x: string; y: number }>, agg: string) {
  if (!data.length) return [];
  // Проверяем, что x — дата
  const dates = data.map(d => new Date(d.x)).filter(d => !isNaN(d.getTime()));
  if (dates.length !== data.length) return data; // не все x — даты
  const minDate = new Date(Math.min(...dates.map(d => d.getTime())));
  const maxDate = new Date(Math.max(...dates.map(d => d.getTime())));
  const days = differenceInDays(maxDate, minDate) + 1;
  const months = differenceInMonths(maxDate, minDate) + 1;
  const years = differenceInYears(maxDate, minDate) + 1;

  let groupKeyFn: (d: Date) => string;
  if (years > 2) {
    groupKeyFn = d => d.getFullYear().toString();
  } else if (months > 3) {
    groupKeyFn = d => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
  } else if (days > 7) {
    groupKeyFn = d => d.toISOString().slice(0, 10);
  } else {
    // до 7 дней — группируем по часам
    groupKeyFn = d => format(d, 'yyyy-MM-dd HH:00');
  }

  const grouped: Record<string, Array<number>> = {};
  data.forEach(({ x, y }) => {
    const date = new Date(x);
    if (!isNaN(date.getTime())) {
      const key = groupKeyFn(date);
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      if (!grouped[key]) grouped[key] = [];
      grouped[key].push(y);
    }
  });
  return Object.entries(grouped).map(([x, ys]) => ({
    x,
    y:
      agg === 'count' || agg === 'sum'
        ? ys.reduce((a, b) => a + b, 0)
        : agg === 'min'
        ? Math.min(...ys)
        : agg === 'max'
        ? Math.max(...ys)
        : ys.reduce((a, b) => a + b, 0) / ys.length
  }));
}

export default function AdminStatsPage() {
  const [x, setX] = useState<AggKey | string>('user_created_at');
  const [y, setY] = useState<AggKey>('count');
  const [agg, setAgg] = useState('count');
  const [data, setData] = useState<Array<{ x: string; y: number }>>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const aggOptions = allowedAgg[y];

  const handleBuild = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await axios.get('/api/admin/stats', {
        params: { x, y, agg },
        headers: { Authorization: `Bearer ${localStorage.getItem('accessToken')}` },
      });
      setData(res.data.points || []);
    } catch (e: any) {
      setError(e?.response?.data?.message || 'Ошибка запроса');
    } finally {
      setLoading(false);
    }
  };

  // Если X — дата, группируем по дню
  const chartData = useMemo(() => {
    let processed = data;
    // Преобразуем время выполнения заказа из наносекунд в часы
    if (y === 'order_completion_time') {
      processed = data.map(d => ({ ...d, y: d.y / 3_600_000_000_000 }));
    }
    if (x.includes('created_at')) {
      const grouped = smartGroupByDate(processed, agg);
      // Сортируем по дате (X)
      return grouped.slice().sort((a, b) => new Date(a.x).getTime() - new Date(b.x).getTime());
    }
    return processed;
  }, [data, x, y, agg]);

  // Формируем данные для chart.js
  const chartJsData = useMemo(() => {
    return {
      labels: chartData.map((d) => d.x),
      datasets: [
        {
          label: allowedY.find((item) => item.value === y)?.label || 'Y',
          data: chartData.map((d) => d.y),
          borderColor: '#8884d8',
          backgroundColor: 'rgba(136,132,216,0.2)',
          tension: 0.3,
        },
      ],
    };
  }, [chartData, y]);

  const chartJsOptions = useMemo(() => {
    const isDate = x.includes('created_at');
    return {
      responsive: true,
      plugins: {
        legend: { display: true },
        title: { display: false },
      },
      scales: {
        x: isDate
          ? {
              type: 'time' as const,
              time: { unit: 'day' as const },
              title: { display: true, text: allowedX.find((item) => item.value === x)?.label || 'X' },
            }
          : {
              title: { display: true, text: allowedX.find((item) => item.value === x)?.label || 'X' },
            },
        y: {
          title: { display: true, text: allowedY.find((item) => item.value === y)?.label || 'Y' },
        },
      },
    };
  }, [x, y]);

  return (
    <div style={{ maxWidth: 900, margin: '0 auto', padding: 24 }}>
      <h2 style={{ marginBottom: 24 }}>Кастомная статистика</h2>
      <div style={{ display: 'flex', gap: 16, marginBottom: 24 }}>
        <Select
          style={{ minWidth: 180 }}
          value={x}
          onChange={setX}
          options={allowedX}
          placeholder="Ось X"
        />
        <Select
          style={{ minWidth: 180 }}
          value={y}
          onChange={val => {
            setY(val as AggKey);
            setAgg(allowedAgg[val as AggKey][0].value);
          }}
          options={allowedY}
          placeholder="Ось Y"
        />
        <Select
          style={{ minWidth: 180 }}
          value={agg}
          onChange={setAgg}
          options={aggOptions}
          placeholder="Агрегация"
        />
        <Button type="primary" onClick={handleBuild} disabled={loading}>
          Построить
        </Button>
      </div>
      {loading && <Spin />}
      {error && <Alert type="error" message={error} style={{ marginBottom: 16 }} />}
      <div style={{ width: '100%', height: 400 }}>
        <Line data={chartJsData} options={chartJsOptions} />
      </div>
    </div>
  );
} 