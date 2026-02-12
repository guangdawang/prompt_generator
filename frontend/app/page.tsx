'use client';

import { useState } from 'react';
import TemplateSelector from '@/components/TemplateSelector';
import VariableInputs from '@/components/VariableInputs';
import TemplateEditor from '@/components/TemplateEditor';
import { Template, templateAPI } from '@/lib/api';

export default function GeneratePage() {
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
  const [variables, setVariables] = useState<Record<string, string>>({});
  const [result, setResult] = useState('');
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [editorOpen, setEditorOpen] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<Template | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const buildInitialValues = (template: Template) => {
    const initialValues: Record<string, string> = {};
    const templateVariables = Array.isArray(template.variables) ? template.variables : [];
    templateVariables.forEach((variable) => {
      initialValues[variable.name] = variable.default_value || '';
    });
    return initialValues;
  };

  const handleTemplateSelect = async (template: Template) => {
    let normalizedTemplate = template;

    if (!Array.isArray(template.variables) || template.variables.length === 0) {
      try {
        const response = await templateAPI.extractVariables(template.content);
        const extracted = (response.variables || []).map((name, index) => ({
          name,
          display_name: name,
          required: true,
          sort_order: index,
        }));

        if (extracted.length > 0) {
          normalizedTemplate = { ...template, variables: extracted };
        }
      } catch (error) {
        console.error('Failed to extract variables:', error);
      }
    }

    setSelectedTemplate(normalizedTemplate);
    setVariables(buildInitialValues(normalizedTemplate));
    setResult('');
  };

  const handleVariableChange = (name: string, value: string) => {
    setVariables(prev => ({ ...prev, [name]: value }));
  };

  const handleGenerate = async () => {
    if (!selectedTemplate) return;

    setLoading(true);
    try {
      const response = await templateAPI.generate({
        template_id: selectedTemplate.id,
        variables,
      });
      setResult(response.result);
    } catch (error) {
      console.error('Failed to generate:', error);
      alert('生成失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    try {
      if (navigator?.clipboard?.writeText) {
        await navigator.clipboard.writeText(result);
      } else {
        const textarea = document.createElement('textarea');
        textarea.value = result;
        textarea.setAttribute('readonly', 'true');
        textarea.style.position = 'absolute';
        textarea.style.left = '-9999px';
        document.body.appendChild(textarea);
        textarea.select();
        const success = document.execCommand('copy');
        document.body.removeChild(textarea);
        if (!success) {
          throw new Error('execCommand copy failed');
        }
      }

      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Failed to copy:', error);
      alert('复制失败，请手动复制');
    }
  };

  const handleBack = () => {
    setSelectedTemplate(null);
    setVariables({});
    setResult('');
  };

  const handleOpenCreate = () => {
    setEditingTemplate(null);
    setEditorOpen(true);
  };

  const handleOpenEdit = () => {
    if (!selectedTemplate) return;
    setEditingTemplate(selectedTemplate);
    setEditorOpen(true);
  };

  const handleEditFromList = (template: Template) => {
    setEditingTemplate(template);
    setEditorOpen(true);
  };

  const handleEditorCancel = () => {
    setEditorOpen(false);
    setEditingTemplate(null);
  };

  const handleEditorSaved = (template: Template) => {
    setEditorOpen(false);
    setEditingTemplate(null);
    setSelectedTemplate(template);
    setVariables(buildInitialValues(template));
    setResult('');
    setRefreshKey((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-gray-900">提示词模板系统</h1>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {editorOpen ? (
          <TemplateEditor
            template={editingTemplate}
            onCancel={handleEditorCancel}
            onSaved={handleEditorSaved}
          />
        ) : !selectedTemplate ? (
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">欢迎使用提示词模板系统</h2>
            <p className="text-gray-600 mb-8">选择一个模板开始生成你的提示词</p>
            <div className="flex items-center justify-between mb-6">
              <span className="text-sm text-gray-500">你也可以创建自己的模板</span>
              <button
                onClick={handleOpenCreate}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                创建新模板
              </button>
            </div>
            <TemplateSelector
              key={refreshKey}
              onSelect={handleTemplateSelect}
              onEdit={handleEditFromList}
            />
          </div>
        ) : (
          <div>
            <div className="flex items-center justify-between mb-6">
              <button
                onClick={handleBack}
                className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
              >
                ← 返回模板列表
              </button>
              <button
                onClick={handleOpenEdit}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
              >
                编辑模板
              </button>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              {/* 左侧：模板信息 */}
              <div className="lg:col-span-1">
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h2 className="text-xl font-semibold mb-2">{selectedTemplate.name}</h2>
                  <p className="text-gray-600 mb-4">{selectedTemplate.description}</p>
                  <div className="border-t pt-4">
                    <h3 className="font-medium text-gray-900 mb-2">模板内容预览</h3>
                    <pre className="bg-gray-50 p-4 rounded text-xs overflow-x-auto whitespace-pre-wrap">
                      {selectedTemplate.content}
                    </pre>
                  </div>
                </div>
              </div>

              {/* 右侧：变量输入和结果 */}
              <div className="lg:col-span-2 space-y-6">
                {/* 变量输入 */}
                <div className="bg-white rounded-lg shadow-md p-6">
                  <h3 className="text-lg font-semibold mb-4">填写变量</h3>
                  <VariableInputs
                    variables={selectedTemplate.variables || []}
                    values={variables}
                    onChange={handleVariableChange}
                  />

                  <div className="mt-6">
                    <button
                      onClick={handleGenerate}
                      disabled={loading}
                      className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
                    >
                      {loading ? '生成中...' : '生成提示词'}
                    </button>
                  </div>
                </div>

                {/* 生成结果 */}
                {result && (
                  <div className="bg-white rounded-lg shadow-md p-6">
                    <div className="flex justify-between items-center mb-4">
                      <h3 className="text-lg font-semibold">生成结果</h3>
                      <button
                        onClick={handleCopy}
                        className="px-4 py-2 text-sm bg-gray-100 rounded hover:bg-gray-200 transition-colors"
                      >
                        {copied ? '已复制！' : '复制'}
                      </button>
                    </div>
                    <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                      <pre className="whitespace-pre-wrap text-sm text-gray-800">
                        {result}
                      </pre>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
