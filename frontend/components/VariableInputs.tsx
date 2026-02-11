'use client';

import { TemplateVariable } from '@/lib/api';

interface VariableInputsProps {
  variables: TemplateVariable[];
  values: Record<string, string>;
  onChange: (name: string, value: string) => void;
}

export default function VariableInputs({ variables, values, onChange }: VariableInputsProps) {
  return (
    <div className="space-y-4">
      {variables.length === 0 ? (
        <p className="text-gray-500 text-center py-4">此模板无需变量</p>
      ) : (
        variables.map((variable) => {
          const inputId = `variable-${variable.name}`;

          return (
            <div key={variable.name} className="bg-gray-50 p-4 rounded-lg">
              <label htmlFor={inputId} className="block text-sm font-medium text-gray-700 mb-2">
                {variable.display_name || variable.name}
                {variable.required && <span className="text-red-500 ml-1">*</span>}
              </label>
              {variable.description && (
                <p className="text-xs text-gray-500 mb-2">{variable.description}</p>
              )}
              {variable.name === 'tone' || variable.name === 'type' || variable.name === 'length' ? (
                <select
                  id={inputId}
                  value={values[variable.name] || variable.default_value || ''}
                  onChange={(e) => onChange(variable.name, e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="">请选择...</option>
                  {variable.name === 'tone' && (
                    <>
                      <option value="正式">正式</option>
                      <option value="友好">友好</option>
                      <option value="简洁">简洁</option>
                      <option value="专业">专业</option>
                    </>
                  )}
                  {variable.name === 'type' && (
                    <>
                      <option value="新闻">新闻</option>
                      <option value="技术">技术</option>
                      <option value="科普">科普</option>
                      <option value="评论">评论</option>
                    </>
                  )}
                  {variable.name === 'length' && (
                    <>
                      <option value="简短">简短</option>
                      <option value="中等">中等</option>
                      <option value="详细">详细</option>
                    </>
                  )}
                </select>
              ) : variable.name === 'content' || variable.name === 'code' ? (
                <textarea
                  id={inputId}
                  value={values[variable.name] || variable.default_value || ''}
                  onChange={(e) => onChange(variable.name, e.target.value)}
                  rows={6}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                  placeholder={variable.default_value || `请输入${variable.display_name}...`}
                />
              ) : (
                <input
                  id={inputId}
                  type="text"
                  value={values[variable.name] || variable.default_value || ''}
                  onChange={(e) => onChange(variable.name, e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder={variable.default_value || `请输入${variable.display_name}...`}
                />
              )}
            </div>
          );
        })
      )}
    </div>
  );
}
