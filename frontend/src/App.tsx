import { ConfigProvider } from 'antd';
import ruRU from 'antd/locale/ru_RU';

function App() {
  return (
    <ConfigProvider locale={ruRU}>
      <div className="text-center">
        <header className="min-h-screen flex flex-col items-center justify-center bg-[#282c34] text-white text-[calc(10px+2vmin)]">
          <h1>Главная страница</h1>
        </header>
      </div>
    </ConfigProvider>
  );
}

export default App;
