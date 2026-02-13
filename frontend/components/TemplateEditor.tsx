'use client';

import { useEffect, useState } from 'react';
import { CreateTemplateRequest, Template, TemplateVariable, templateAPI } from '@/lib/api';

interface TemplateEditorProps {
  template?: Template | null;
  onCancel: () => void;
  onSaved: (template: Template) => void;
}

const createEmptyVariable = (): TemplateVariable => ({
  name: '',
  display_name: '',
  description: '',
  default_value: '',
  required: true,
  sort_order: 0,
});

export default function TemplateEditor({ template, onCancel, onSaved }: TemplateEditorProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [isPublic, setIsPublic] = useState(false);
  const [content, setContent] = useState('');
  const [variables, setVariables] = useState<TemplateVariable[]>([]);
  const [saving, setSaving] = useState(false);
  const [extracting, setExtracting] = useState(false);

  useEffect(() => {
    setName(template?.name || '');
    setDescription(template?.description || '');
    setCategory(template?.category || '');
    setIsPublic(Boolean(template?.is_public));
    setContent(template?.content || '');
    setVariables(template?.variables || []);
  }, [template]);

  const handleAddVariable = () => {
    setVariables((prev) => [...prev, createEmptyVariable()]);
  };

  const handleRemoveVariable = (index: number) => {
    setVariables((prev) => prev.filter((_, i) => i !== index));
  };

  const handleVariableChange = (index: number, updates: Partial<TemplateVariable>) => {
    setVariables((prev) => {
      const next = [...prev];
      next[index] = { ...next[index], ...updates };
      return next;
    });
  };

  const handleExtractVariables = async () => {
    if (!content.trim()) return;

    setExtracting(true);
    try {
      const response = await templateAPI.extractVariables(content);
      const extracted = response.variables || [];

      setVariables((prev) => {
        const existing = new Map(prev.map((variable) => [variable.name, variable]));
        return extracted.map((name, index) => {
          const current = existing.get(name);
          return {
            name,
            display_name: current?.display_name || name,
            description: current?.description || '',
            default_value: current?.default_value || '',
            required: current?.required ?? true,
            sort_order: index,
          };
        });
      });
    } catch (error) {
      console.error('Failed to extract variables:', error);
      alert('变量提取失败，请检查模板内容');
    } finally {
      setExtracting(false);
    }
  };

  const handleSave = async () => {
    if (!name.trim() || !content.trim()) {
      alert('请填写模板名称和模板内容');
      return;
    }

    const sanitizedVariables = variables
      .filter((variable) => variable.name.trim() !== '')
      .map((variable, index) => ({
        ...variable,
        name: variable.name.trim(),
        display_name: (variable.display_name || variable.name).trim(),
        sort_order: index,
        required: variable.required ?? true,
      }));

    // 客户端校验变量名，确保符合后端要求：^[a-zA-Z_][a-zA-Z0-9_]*$
    const namePattern = /^[a-zA-Z_][a-zA-Z0-9_]*$/;
    const invalid = sanitizedVariables.filter((v) => !namePattern.test(v.name));
    if (invalid.length > 0) {
      const names = invalid.map((v) => v.name || '(empty)').join(', ');
      alert(`变量名格式错误：${names}\n变量名只能包含字母、数字和下划线，且不能以数字开头。`);
      return;
    }

    const payload: CreateTemplateRequest = {
      name: name.trim(),
      description: description.trim(),
      content,
      variables: sanitizedVariables,
      category: category.trim(),
      is_public: isPublic,
    };

    setSaving(true);
    try {
      const saved = template
        ? await templateAPI.updateTemplate(template.id, payload)
        : await templateAPI.createTemplate(payload);
      onSaved(saved);
    } catch (error) {
      console.error('Failed to save template:', error);
      alert('保存失败，请重试');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">
          {template ? '编辑模板' : '创建模板'}
        </h2>
        <div className="flex items-center gap-2">
          <button
            onClick={onCancel}
            className="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200"
          >
            取消
          </button>
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {saving ? '保存中...' : '保存模板'}
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">模板名称</label>
          <input
            type="text"
            value={name}
            onChange={(event) => setName(event.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="请输入模板名称"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">分类</label>
          <input
            type="text"
            value={category}
            onChange={(event) => setCategory(event.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="例如：写作、总结"
          />
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">模板描述</label>
        <textarea
          value={description}
          onChange={(event) => setDescription(event.target.value)}
          rows={3}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="简要描述模板用途"
        />
      </div>

      <div>
        <div className="flex items-center justify-between mb-2">
          <label className="block text-sm font-medium text-gray-700">模板内容</label>
          <button
            onClick={handleExtractVariables}
            disabled={extracting || !content.trim()}
            className="px-3 py-1 text-sm bg-gray-100 rounded hover:bg-gray-200 disabled:opacity-50"
          >
            {extracting ? '提取中...' : '从内容提取变量'}
          </button>
        </div>
        <textarea
          value={content}
          onChange={(event) => setContent(event.target.value)}
          rows={8}
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 font-mono text-sm"
          placeholder="在这里输入模板内容，例如：请总结 {{topic}} 的核心观点"
        />
        <label className="flex items-center gap-2 mt-3 text-sm text-gray-600">
          <input
            type="checkbox"
            checked={isPublic}
            onChange={(event) => setIsPublic(event.target.checked)}
            className="rounded border-gray-300"
          />
          公开模板
        </label>
      </div>

      <div>
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-lg font-semibold">变量设置</h3>
          <button
            onClick={handleAddVariable}
            className="px-3 py-1 text-sm bg-blue-50 text-blue-700 rounded hover:bg-blue-100"
          >
            添加变量
          </button>
        </div>
        {variables.length === 0 ? (
          <p className="text-gray-500 text-sm">暂无变量，点击上方按钮添加或从内容提取。</p>
        ) : (
          <div className="space-y-3">
            {variables.map((variable, index) => (
              <div key={`${variable.name}-${index}`} className="border border-gray-200 rounded-lg p-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div>
                    <label className="block text-xs text-gray-600 mb-1">变量名</label>
                    <input
                      type="text"
                      value={variable.name}
                      onChange={(event) => handleVariableChange(index, { name: event.target.value })}
                      className="w-full px-2 py-1 border border-gray-300 rounded"
                      placeholder="例如：topic"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-gray-600 mb-1">显示名称</label>
                    <input
                      type="text"
                      value={variable.display_name}
                      onChange={(event) => handleVariableChange(index, { display_name: event.target.value })}
                      className="w-full px-2 py-1 border border-gray-300 rounded"
                      placeholder="例如：主题"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-gray-600 mb-1">默认值</label>
                    <input
                      type="text"
                      value={variable.default_value || ''}
                      onChange={(event) => handleVariableChange(index, { default_value: event.target.value })}
                      className="w-full px-2 py-1 border border-gray-300 rounded"
                      placeholder="可选"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-gray-600 mb-1">说明</label>
                    <input
                      type="text"
                      value={variable.description || ''}
                      onChange={(event) => handleVariableChange(index, { description: event.target.value })}
                      className="w-full px-2 py-1 border border-gray-300 rounded"
                      placeholder="可选"
                    />
                  </div>
                </div>
                <div className="flex items-center justify-between mt-3">
                  <label className="flex items-center gap-2 text-xs text-gray-600">
                    <input
                      type="checkbox"
                      checked={variable.required}
                      onChange={(event) => handleVariableChange(index, { required: event.target.checked })}
                      className="rounded border-gray-300"
                    />
                    必填
                  </label>
                  <button
                    onClick={() => handleRemoveVariable(index)}
                    className="text-xs text-red-600 hover:text-red-700"
                  >
                    删除变量
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
