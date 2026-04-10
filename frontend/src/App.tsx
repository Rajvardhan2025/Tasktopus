import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Landing } from './pages/Landing';
import { ProjectList } from './pages/ProjectList';
import { ProjectBoard } from './pages/ProjectBoard';
import Login from './pages/Login';
import { Layout } from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import { Toaster } from './components/ui/toaster';
import { authService } from './lib/auth';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

function App() {
  const isAuthenticated = authService.isAuthenticated();

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route 
            path="/" 
            element={isAuthenticated ? <Navigate to="/projects" replace /> : <Login />} 
          />
          <Route 
            path="/login" 
            element={isAuthenticated ? <Navigate to="/projects" replace /> : <Login />} 
          />
          <Route element={<Layout />}>
            <Route 
              path="projects" 
              element={
                <ProtectedRoute>
                  <ProjectList />
                </ProtectedRoute>
              } 
            />
            <Route 
              path="projects/:projectId" 
              element={
                <ProtectedRoute>
                  <ProjectBoard />
                </ProtectedRoute>
              } 
            />
          </Route>
        </Routes>
      </BrowserRouter>
      <Toaster />
    </QueryClientProvider>
  );
}

export default App;
