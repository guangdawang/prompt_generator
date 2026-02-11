import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export interface Template {
  id: string;
  user_id: string;
  name: string;
  description: string;
  content: string;
  variables: TemplateVariable[];
  category: string;
  is_public: boolean;
  usage_count: number;
  created_at: string;
  updated_at: string;
}

type RawTemplate = Omit<Template, 'variables'> & {
  variables?: unknown;
};

export interface TemplateVariable {
  id?: string;
  name: string;
  display_name: string;
  description?: string;
  default_value?: string;
  required: boolean;
  sort_order?: number;
}

export interface GenerateRequest {
  template_id: string;
  variables: Record<string, string>;
}

export interface GenerateResponse {
  result: string;
  prompt: string;
}

export interface CreateTemplateRequest {
  name: string;
  description?: string;
  content: string;
  variables?: TemplateVariable[];
  category?: string;
  is_public?: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total: number;
}

export const templateAPI = {
  // 健康检查
  health: async () => {
    return api.get('/health');
  },

  // 获取模板列表
  getTemplates: async (category?: string, page = 1, pageSize = 20): Promise<PaginatedResponse<Template>> => {
    const params: any = { page, page_size: pageSize };
    if (category) params.category = category;
    const response = await api.get('/templates', { params }) as PaginatedResponse<RawTemplate>;
    return {
      ...response,
      data: normalizeTemplateList(response.data || []),
    };
  },

  // 获取公开模板
  getPublicTemplates: async (category?: string, page = 1, pageSize = 20): Promise<PaginatedResponse<Template>> => {
    const params: any = { page, page_size: pageSize };
    if (category) params.category = category;
    const response = await api.get('/templates/public', { params }) as PaginatedResponse<RawTemplate>;
    return {
      ...response,
      data: normalizeTemplateList(response.data || []),
    };
  },

  // 获取单个模板
  getTemplate: async (id: string): Promise<Template> => {
    const response = await api.get(`/templates/${id}`) as RawTemplate;
    return normalizeTemplate(response);
  },

  // 创建模板
  createTemplate: async (template: CreateTemplateRequest): Promise<Template> => {
    const response = await api.post('/templates', template) as RawTemplate;
    return normalizeTemplate(response);
  },

  // 更新模板
  updateTemplate: async (id: string, template: Partial<CreateTemplateRequest>): Promise<Template> => {
    const response = await api.put(`/templates/${id}`, template) as RawTemplate;
    return normalizeTemplate(response);
  },

  // 删除模板
  deleteTemplate: async (id: string): Promise<void> => {
    return api.delete(`/templates/${id}`);
  },

  // 生成提示词
  generate: async (data: GenerateRequest): Promise<GenerateResponse> => {
    return api.post('/generate', data);
  },

  // 提取变量
  extractVariables: async (content: string): Promise<{ variables: string[] }> => {
    return api.post('/generate/extract-variables', { content });
  },
};

const normalizeTemplateList = (templates: RawTemplate[]): Template[] => {
  return templates.map(normalizeTemplate);
};

const normalizeTemplate = (template: RawTemplate): Template => {
  return {
    ...template,
    variables: normalizeVariables(template.variables),
  };
};

const normalizeVariables = (variables: unknown): TemplateVariable[] => {
  if (Array.isArray(variables)) {
    return variables as TemplateVariable[];
  }
  if (typeof variables === 'string') {
    try {
      const parsed = JSON.parse(variables);
      return Array.isArray(parsed) ? (parsed as TemplateVariable[]) : [];
    } catch {
      return [];
    }
  }
  return [];
};

export default templateAPI;
