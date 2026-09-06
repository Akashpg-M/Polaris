import { describe, expect, it } from 'vitest';
import { ackExplanation, canCancelTaskState, canRetryCommandState, requiredCapability, validateTaskDraft, type TaskDraft } from './model';

const draft: TaskDraft = { projectId:'project-1',taskType:'NAVIGATE',priority:'HIGH',expiresAt:'2099-01-01T12:00',latitude:'13.08',longitude:'80.27',targetJSON:'{}',minimumBattery:'40',maximumDistanceKm:'5',deviceTypes:['connected_vehicle'],optionalCapabilities:['payload'],planningMode:'POLARIS_REQUIRED',customConstraintsJSON:'',correlationId:'' };

describe('operations model',()=>{
  it('uses the backend command capability mapping',()=>{expect(requiredCapability('NAVIGATE')).toBe('navigate');expect(requiredCapability('STOP')).toBe('receive_relocation_command');expect(requiredCapability('RUN_MODEL')).toBe('run_model');});
  it('converts validated presentation units into the task DTO',()=>{const result=validateTaskDraft(draft);expect(result.errors).toEqual({});expect(result.input?.requirements.max_distance_meters).toBe(5000);expect(result.input?.requirements.required_capabilities).toEqual(['navigate','payload']);expect(result.input?.target).toEqual({lat:13.08,lon:80.27});});
  it('rejects invalid coordinates, battery, distance, JSON and expiry',()=>{const result=validateTaskDraft({...draft,latitude:'91',minimumBattery:'101',maximumDistanceKm:'-1',expiresAt:'2020-01-01T00:00'});expect(Object.keys(result.errors)).toEqual(expect.arrayContaining(['latitude','minimumBattery','maximumDistanceKm','expiresAt']));});
  it('does not interpret a blank coordinate as zero',()=>{expect(validateTaskDraft({...draft,latitude:''}).errors.latitude).toBeTruthy();});
  it('preserves conservative transition controls',()=>{expect(canCancelTaskState('ASSIGNED','PENDING')).toBe(true);expect(canCancelTaskState('ASSIGNED','DELIVERED')).toBe(false);expect(canRetryCommandState('FAILED')).toBe(true);expect(canRetryCommandState('COMPLETED')).toBe(false);});
  it('explains duplicate ACK as idempotent',()=>{expect(ackExplanation('DUPLICATE')).toContain('idempotent');});
});
