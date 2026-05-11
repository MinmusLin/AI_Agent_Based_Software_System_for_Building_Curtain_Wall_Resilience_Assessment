import type { DragEvent, ReactElement } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';

import { AppLayout } from './components/AppLayout';
import DashboardPage from './pages/DashboardPage';
import ForgetPasswordPage from './pages/ForgetPasswordPage';
import LoginPage from './pages/LoginPage';
import NotFoundPage from './pages/NotFoundPage';
import ProjectDetailPage from './pages/ProjectDetailPage';
import ProjectsPage from './pages/ProjectsPage';
import RegisterPage from './pages/RegisterPage';
import { GuestRoute, ProtectedRoute } from './routes/RouteGuards';

function preventNativeImageDrag(event: DragEvent<HTMLDivElement>): void {
  if (event.target instanceof HTMLImageElement) {
    event.preventDefault();
  }
}

export default function App(): ReactElement {
  return (
    <div className="contents" onDragStartCapture={preventNativeImageDrag}>
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
            <Route element={<ProjectDetailPage />} path="/projects/:projectId" />
            <Route element={<ProjectDetailPage />} path="/projects/:projectId/:stage" />
          </Route>
        </Route>
        <Route element={<NotFoundPage />} path="*" />
      </Routes>
    </div>
  );
}
