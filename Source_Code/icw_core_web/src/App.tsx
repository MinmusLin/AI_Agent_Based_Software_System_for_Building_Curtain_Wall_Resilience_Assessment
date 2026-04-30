import type { ReactElement } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';

import { AppLayout } from './components/AppLayout';
import DashboardPage from './pages/DashboardPage';
import ForgetPasswordPage from './pages/ForgetPasswordPage';
import LoginPage from './pages/LoginPage';
import ProjectsPage from './pages/ProjectsPage';
import RegisterPage from './pages/RegisterPage';
import { GuestRoute, ProtectedRoute } from './routes/RouteGuards';

export default function App(): ReactElement {
  return (
    <Routes>
      <Route element={<Navigate replace to="/dashboard" />} path="/" />
      <Route element={<GuestRoute />}>
        <Route element={<LoginPage />} path="/login" />
        <Route element={<RegisterPage />} path="/register" />
        <Route element={<ForgetPasswordPage />} path="/forget-password" />
      </Route>
      <Route element={<ProtectedRoute />}>
        <Route element={<AppLayout />}>
          <Route element={<DashboardPage />} path="/dashboard" />
          <Route element={<ProjectsPage />} path="/projects" />
        </Route>
      </Route>
      <Route element={<Navigate replace to="/dashboard" />} path="*" />
    </Routes>
  );
}
